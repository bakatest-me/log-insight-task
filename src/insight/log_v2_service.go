package insight

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log-insign-task/domain"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var location, _ = time.LoadLocation("Asia/Bangkok")

type logV2Service struct{}

func NewLogV2Service() *logV2Service {
	return &logV2Service{}
}

func (m *logV2Service) Run(pathFile string) (*domain.SumaryLog, error) {
	// Open the log file
	file, err := os.Open(pathFile)
	if err != nil {
		return nil, errors.New("Error opening file")
	}
	defer file.Close()

	lineCh := make(chan domain.Log)
	lineSummaryCh := make(chan domain.SumaryLog)
	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			m.processLog(lineCh, lineSummaryCh)
			wg.Done()
		}()
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		var logLine domain.Log
		if err := json.Unmarshal(line, &logLine); err != nil {
			fmt.Println("Error parsing log line:", err)
			continue
		}
		lineCh <- logLine
	}
	close(lineCh)

	if err := scanner.Err(); err != nil {
		return nil, errors.New("Error reading file")
	}

	go func() {
		wg.Wait()
		close(lineSummaryCh)
	}()

	return m.summary(lineSummaryCh), nil
}

func (m *logV2Service) summary(lineSummaryCh <-chan domain.SumaryLog) *domain.SumaryLog {
	resp := domain.SumaryLog{
		LogLevelFrequency: map[string]int{},
		HttpCodeFrequency: map[int]int{},
		FrequencyByTime:   map[int]int{},
		UriFrequency:      map[string]int{},
	}

	var max, min int
	for v := range lineSummaryCh {
		if max <= 0 {
			max = v.Max
		}
		if min <= 0 {
			min = v.Min
		}

		resp.TotalLog += v.TotalLog
		resp.LogLevelFrequency = mergeMap(resp.LogLevelFrequency, v.LogLevelFrequency)
		resp.HttpCodeFrequency = mergeMapInt(resp.HttpCodeFrequency, v.HttpCodeFrequency)
		resp.FrequencyByTime = mergeMapInt(resp.FrequencyByTime, v.FrequencyByTime)
		resp.UriFrequency = mergeMap(resp.UriFrequency, v.UriFrequency)
		resp.Max = m.calculateMax(resp.Max, v.Max)
		resp.Min = m.calculateMin(resp.Min, v.Min)
		resp.SumLatency += v.SumLatency
		resp.TotalLogLatency += v.TotalLogLatency
		resp.TotalLongLatency += v.TotalLongLatency
	}
	return &resp

}

func (m *logV2Service) processLog(lineCh chan domain.Log, summaryCh chan<- domain.SumaryLog) {
	summary := domain.SumaryLog{
		LogLevelFrequency: make(map[string]int),
		HttpCodeFrequency: make(map[int]int),
		UriFrequency:      make(map[string]int),
		FrequencyByTime:   make(map[int]int),
	}

	longLatency := 500

	for logLine := range lineCh {
		summary.TotalLog++
		err := m.updateFrequencies(&summary, &logLine)
		if err != nil {
			// don't count status zero
			continue
		}
		ms, err := m.getLatency(logLine.Latency)
		if err == nil {
			summary.TotalLongLatency = m.countLongLatency(ms, summary.TotalLongLatency, longLatency)
			summary.Max = m.calculateMax(summary.Max, ms)
			summary.Min = m.calculateMin(summary.Min, ms)

			summary.TotalLogLatency++
			summary.SumLatency += ms
		}

		hour, err := m.extractHour(logLine.Ts, "Asia/Bangkok")
		if err != nil {
			fmt.Println("Error parsing log line on ts field:", logLine, err)
			continue
		}
		summary.FrequencyByTime[hour]++
	}
	summaryCh <- summary
}

func (m *logV2Service) updateFrequencies(summary *domain.SumaryLog, logLine *domain.Log) error {
	summary.LogLevelFrequency[logLine.Level]++

	if logLine.Status <= 0 {
		return errors.New("status is zero")
	}

	summary.HttpCodeFrequency[logLine.Status]++
	path := m.extractPath(logLine.URI)
	summary.UriFrequency[path]++

	return nil
}

func (m *logV2Service) extractHour(ts string, tz string) (int, error) {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return 0, err
	}

	t = t.In(location)
	return t.Hour(), nil
}

func (m *logV2Service) extractPath(uri string) string {
	splitPath := strings.Split(uri, "/")
	path := strings.Join(splitPath[0:4], "/")
	return path
}

func (m *logV2Service) countLongLatency(ms int, totalLongLatency int, lognLatency int) int {
	if ms >= lognLatency {
		totalLongLatency++
	}
	return totalLongLatency
}

func (m *logV2Service) calculateMax(max int, ms int) int {
	if ms > max {
		max = ms
	}
	return max
}

func (m *logV2Service) calculateMin(min int, ms int) int {
	if ms < min || min == 0 {
		min = ms
	}
	return min
}

func (m *logV2Service) getLatency(latency string) (ms int, err error) {
	duration, err := time.ParseDuration(latency)
	if err != nil {
		return 0, err
	}

	ms = int(duration / time.Millisecond)
	if ms == 0 {
		return 0, errors.New("latency is zero")
	}
	return ms, nil
}

func (m *logV2Service) Print(resp *domain.SumaryLog) {
	fmt.Println("Total log:", resp.TotalLog)
	fmt.Println("Min latency:", resp.Min)
	fmt.Println("Max latency:", resp.Max)
	fmt.Println("Avg latency:", resp.GetAvg())
	fmt.Println("Total log latency (more then equal 500ms): ", resp.TotalLongLatency)

	fmt.Println()
	fmt.Println("Log level frequency:")
	Print(resp.LogLevelFrequency)

	fmt.Println()
	fmt.Println("Http code frequency:")
	Print(sortByKeyInt(resp.HttpCodeFrequency))

	fmt.Println()
	fmt.Println("Top 5 URI")
	for i, v := range sortTopRank(resp.UriFrequency) {
		fmt.Printf("%v. %v: %v\n", i+1, v.Key, v.Value)
		if i == 4 {
			break
		}
	}

	fmt.Println()
	fmt.Println("Frequency by time in 24 hour:")
	for _, v := range sortByKeyInt(resp.FrequencyByTime) {
		fmt.Printf("Hour at %v: %v\n", v.Key, v.Value)
	}

}
