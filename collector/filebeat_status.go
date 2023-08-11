package collector

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"regexp"
	"strings"
)

type Harvester struct {
	OpenFiles int `json:"open_files"`
	Running   int `json:"running"`
}

type Filebeat struct {
	Harvester Harvester `json:"harvester"`
}

type FilebeatStatusCollector struct {
	FilebeatStatus *prometheus.Desc
}

func init() {
	registerCollector("filebeat", NewFilebeatStatusCollector, "Filebeat")
}

func NewFilebeatStatusCollector() (Collector, error) {
	const subsystem = "filebeat"

	return &FilebeatStatusCollector{
		FilebeatStatus: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "openfiles"),
			"Filebeat monitoring log harvester openfiles running.",
			nil,
			nil,
		),
	}, nil
}

func (c *FilebeatStatusCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	openfilesValue := checkFilebeatStatus()
	ch <- prometheus.MustNewConstMetric(
		c.FilebeatStatus,
		prometheus.CounterValue,
		openfilesValue,
	)
	return nil
}

//check filebeat openfiles returns type float64
func checkFilebeatStatus() float64 {
	var openFiles float64 = 0.0
	logPath := "C:\\ProgramData\\Elastic\\Beats\\filebeat\\logs\\filebeat"
	file, err := os.Open(logPath)
	if err != nil {
		fmt.Println("Filebeat_openErr:", err)
		return openFiles
	}
	defer file.Close()
	buf := make([]byte, 204800)
	stat, err := os.Stat(logPath)
	if err != nil {
		fmt.Println("Filebeat_statErr:", err)
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

	re := regexp.MustCompile(`"open_files":\d+,"running":\d+`)
	result := re.FindString(lastLine)
	var out string
	if result == "" {
		out = `{"harvester":{"open_files":0,"running":0}}`
	} else {
		out = `{"harvester":{` + result + `}}`
	}
	var f Filebeat
	err = json.Unmarshal([]byte(out), &f)
	if err != nil {
		fmt.Println("Filebeat_unmarshalErr:", err)
	}
	outInt := f.Harvester.Running
	outFloat := float64(outInt)
	openFiles = outFloat
	return openFiles
}
