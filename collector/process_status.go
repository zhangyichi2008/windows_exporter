package collector

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	allProcess = []string{
		"filebeat",
		"winlogbeat",
	}
)

func init() {
	registerCollector("process_status", NewProcessStatusCollector, "Process_Status")
}

type ProcessStatusCollector struct {
	ProcessStatus *prometheus.Desc
}

func NewProcessStatusCollector() (Collector, error) {
	const subsystem = "process_status"

	return &ProcessStatusCollector{
		ProcessStatus: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "up"),
			"Value is 1 if the process is up, 0 otherwise.",
			[]string{"process_name"},
			nil,
		),
	}, nil
}

func (c *ProcessStatusCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	for _, process := range allProcess {
		upValue := 0.0
		alive := checkProcess(process)
		if alive {
			upValue = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.ProcessStatus,
			prometheus.CounterValue,
			upValue,
			process,
		)

	}
	return nil
}

func checkProcess(name string) bool {
	c, b := exec.Command("tasklist"), new(bytes.Buffer)
	c.Stdout = b
	c.Run()
	s := bufio.NewScanner(b)
	for s.Scan() {
		if strings.Contains(s.Text(), name+".exe") {
			return true
		}
	}
	return false
}
