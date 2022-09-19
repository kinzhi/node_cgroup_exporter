package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	cgroupMemSubsystem = "cgroupMem"
)

type cgroupmemCollector struct {
	logger log.Logger
}

func init() {
	registerCollector("cgroupusagemem", defaultEnabled, NewCgroupmemCollector)
}

// NewCgroupmemCollector returns a new Collector exposing memory stats.
func NewCgroupmemCollector(logger log.Logger) (Collector, error) {
	return &cgroupmemCollector{logger}, nil
}

// Update calls (*cgroupmemCollector).getCgroupUsageMem to get the platform specific
// memory metrics.
func (c *cgroupmemCollector) Update(ch chan<- prometheus.Metric) error {
	// var metricType prometheus.ValueType
	cgroupUsageMem, err := c.getCgroupUsageMem()
	if err != nil {
		return fmt.Errorf("couldn't get cgroupmemusage: %w", err)
	}
	level.Debug(c.logger).Log("msg", "Set cgroup_mem", "cgroupMem", cgroupUsageMem)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cgroupMemSubsystem, "usage"),
			fmt.Sprintf("Memory information field %s.", cgroupUsageMem),
			nil, nil,
		), prometheus.GaugeValue, cgroupUsageMem,
	)
	return nil
}

// 	ch <- prometheus.MustNewConstMetric(h.memDesc, prometheus.GaugeValue, 80, h.labelVaues...)
