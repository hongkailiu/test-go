package cincinnati

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricsPrefix = "cincinnati_go"
)

var (
	metrics = struct {
		tagScraped *prometheus.CounterVec
	}{
		tagScraped: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: MetricsPrefix + "_tag_scraped_total",
			Help: "Total number of tags scraped",
		}, []string{"at"}),
	}
)

func init() {
	prometheus.MustRegister(metrics.tagScraped)
}
