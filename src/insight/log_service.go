package insight

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log-insign-task/domain"
	"os"
	"strings"
	"time"
)

type logService struct{}

func NewLogService() *logService {
	return &logService{}
}

func (m *logService) Run(pathFile string) (*domain.SumaryLog, error) {
	// Open the log file
	file, err := os.Open(pathFile)
	if err != nil {
		return nil, errors.New("Error opening file")
	}
	defer file.Close()

	logLevelFrequency := make(map[string]int)
	httpCodeFrequency := make(map[int]int)
	uriFrequency := make(map[string]int)
	frequencyByTime := make(map[int]int)

	lineCount := 0
	var min, max int
	var sumLatency, totalLogLetency int
	totalLongLatency := 0
	lognLatency := 500

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
		line := scanner.Bytes()
		var logLine domain.Log
		if err := json.Unmarshal(line, &logLine); err != nil {
			fmt.Println("Error parsing log line:", err)
			continue
		}

		logLevelFrequency[logLine.Level]++

		if logLine.Status <= 0 {
			continue
		}

		httpCodeFrequency[logLine.Status]++

		ms, err := m.getLatency(logLine.Latency)
		if err == nil {
			totalLongLatency = m.countLongLatency(ms, totalLongLatency, lognLatency)
			max = m.calculateMax(max, ms)
			min = m.calculateMin(min, ms)

			totalLogLetency++
			sumLatency += ms
		}

		path := m.extractPath(logLine.URI)
		uriFrequency[path]++

		hour, err := m.extractHour(logLine.Ts, "Asia/Bangkok")
		if err != nil {
			fmt.Println("Error parsing log line on ts field:", logLine, err)
			continue
		}
		frequencyByTime[hour]++
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.New("Error reading file")
	}

	return &domain.SumaryLog{
		TotalLog:          lineCount,
		LogLevelFrequency: logLevelFrequency,
		HttpCodeFrequency: httpCodeFrequency,
		FrequencyByTime:   frequencyByTime,
		UriFrequency:      uriFrequency,
		Min:               min,
		Max:               max,
		Avg:               0,
		SumLatency:        sumLatency,
		TotalLogLatency:   totalLogLetency,
		TotalLongLatency:  totalLongLatency,
	}, nil
}

func (m *logService) extractHour(ts string, tz string) (int, error) {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return 0, err
	}

	location, err := time.LoadLocation(tz)
	if err != nil {
		return 0, err
	}

	t = t.In(location)
	return t.Hour(), nil
}

func (m *logService) extractPath(uri string) string {
	splitPath := strings.Split(uri, "/")
	path := strings.Join(splitPath[0:4], "/")
	return path
}

func (m *logService) countLongLatency(ms int, totalLongLatency int, lognLatency int) int {
	if ms >= lognLatency {
		totalLongLatency++
	}
	return totalLongLatency
}

func (m *logService) calculateMax(max int, ms int) int {
	if ms > max {
		max = ms
	}
	return max
}

func (m *logService) calculateMin(min int, ms int) int {
	if ms < min || min == 0 {
		min = ms
	}
	return min
}

func (m *logService) getLatency(latency string) (ms int, err error) {
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

func (m *logService) Print(resp *domain.SumaryLog) {
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
	Print(resp.HttpCodeFrequency)

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
