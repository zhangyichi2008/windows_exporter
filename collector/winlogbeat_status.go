package collector

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"regexp"
	"strings"
)

type Published struct {
	Total int `json:"total"`
}

type Winlogbeat struct {
	Published Published `json:"published_events"`
}

type WinlogbeatStatusCollector struct {
	WinlogbeatStatus *prometheus.Desc
}

func init() {
	registerCollector("winlogbeat", NewWinlogbeatStatusCollector, "Winlogbeat")
}

func NewWinlogbeatStatusCollector() (Collector, error) {
	const subsystem = "winlogbeat"

	return &WinlogbeatStatusCollector{
		WinlogbeatStatus: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "published_events"),
			"Winlogbeat monitoring log total published events.",
			nil,
			nil,
		),
	}, nil
}

func (c *WinlogbeatStatusCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	publishedEventsValue := checkWinlogbeatStatus()
	ch <- prometheus.MustNewConstMetric(
		c.WinlogbeatStatus,
		prometheus.CounterValue,
		publishedEventsValue,
	)
	return nil
}

func checkWinlogbeatStatus() float64 {
	var publishedEvents float64
	logPath := "C:\\winlogbeat-7.12.0-windows-x86_64\\logs\\winlogbeat"
	file, err := os.Open(logPath)
	if err != nil {
		fmt.Println("Winlogbeat_openErr:", err)
	}
	defer file.Close()

	buf := make([]byte, 204800)
	stat, err := os.Stat(logPath)
	if err != nil {
		fmt.Println("Winlogbeat_statErr:", err)
	}
	start := stat.Size() - 204800
	if start < 0 {
		start = 0
	}
	n, _ := file.ReadAt(buf, start)
	lines := strings.Split(string(buf[:n]), "\n")
	var lastLine string
	for _, line := range lines {
		if strings.Contains(line, "monitoring") {
			lastLine = line
		}
	}
	re := regexp.MustCompile(`"total":\d+}}}}`)
	result := strings.Split(re.FindString(lastLine), "}")[0]
	var out string
	if result == "" {
		out = `{"published_events":{"total":0}}`
	} else {
		out = `{"published_events":{` + result + `}}`
	}
	var f Winlogbeat
	err = json.Unmarshal([]byte(out), &f)
	if err != nil {
		fmt.Println("Winlogbeat_unmarshalErr:", err)
	}
	outInt := f.Published.Total
	outFloat := float64(outInt)
	publishedEvents = outFloat
	return publishedEvents
}
