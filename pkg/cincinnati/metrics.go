package cincinnati

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var (
	tagScrapedMetric = &ginprometheus.Metric{
		ID:          "tag_scraped_total",
		Name:        "tag_scraped_total",
		Description: "Total number of tags scraped",
		Type:        "counter_vec",
		Args:        []string{"at"},
	}
)

type metrics struct {
	tagScrapedCounterVec *prometheus.CounterVec

	p *ginprometheus.Prometheus
}

func (m *metrics) tagScraped() *prometheus.CounterVec {
	if m.tagScrapedCounterVec != nil {
		return m.tagScrapedCounterVec
	}
	for _, met := range m.p.MetricsList {
		if met.ID == "tag_scraped_total" {
			m.tagScrapedCounterVec = met.MetricCollector.(*prometheus.CounterVec)
			return m.tagScrapedCounterVec
		}
	}
	panic("failed to find metrics with id tag_scraped_total")
}

func NewMetrics(address string, g *gin.Engine) *metrics {

	p := ginprometheus.NewWithConfig(ginprometheus.Config{
		Subsystem:          "cincinnati",
		DisableBodyReading: true,
		MetricsList:        []*ginprometheus.Metric{tagScrapedMetric},
	})
	p.SetListenAddress(address)
	p.Use(g)

	return &metrics{p: p}
}
