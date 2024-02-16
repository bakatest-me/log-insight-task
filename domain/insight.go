package domain

type Log struct {
	URI          string `json:"URI"`
	Caller       string `json:"caller"`
	Header       string `json:"header"`
	Host         string `json:"host"`
	Latency      string `json:"latency"`
	Level        string `json:"level"`
	Method       string `json:"method"`
	Msg          string `json:"msg"`
	RealIP       string `json:"real_ip"`
	RequestID    string `json:"request_id"`
	ResponseSize string `json:"response_size"`
	Status       int    `json:"status"`
	Ts           string `json:"ts"`
}

type LogLevel string

var (
	LevelInfo  LogLevel = "info"
	LevelError LogLevel = "error"
	LevelWarn  LogLevel = "warn"
	LevelDebug LogLevel = "debug"
)

type SumaryLog struct {
	TotalLog          int
	LogLevelFrequency map[string]int
	HttpCodeFrequency map[int]int
	FrequencyByTime   map[int]int
	UriFrequency      map[string]int
	Min               int
	Max               int
	Avg               int
	SumLatency        int
	TotalLogLatency   int
	TotalLongLatency  int
}

func (m SumaryLog) GetAvg() int {
	return m.SumLatency / m.TotalLogLatency
}
