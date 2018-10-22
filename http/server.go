package http

import (
	"github.com/gin-gonic/gin"
	"github.com/hongkailiu/test-go/random"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"path"},
	)

	randomNumber = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "random_number",
			Help: "the value of random number.",
		}, []string{"key"},

	)

	//https://github.com/kubernetes/kubernetes/blob/master/pkg/volume/util/metrics.go
	storageOperationMetric = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "storage_operation_duration_seconds",
			Help:    "Storage operation duration",
			Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10, 15, 25, 50},
		},
		[]string{"volume_plugin", "operation_name"},
	)
)

func PrometheusRegister() {
	log.WithFields(log.Fields{
		"name": "httpRequestsTotal",
	}).Info("prometheus register")
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(randomNumber)
	randomNumber.With(prometheus.Labels{"key": "value"}).Set(0)
	prometheus.MustRegister(storageOperationMetric)
}

func PrometheusLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithFields(log.Fields{
			"c.Request.URL.Path": c.Request.URL.Path,
		}).Debug("prometheus logger detected path visited")
		httpRequestsTotal.With(prometheus.Labels{"path": c.Request.URL.Path}).Inc()
	}
}

func Run() {
	PrometheusRegister()

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	gin.Logger()
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.Use(PrometheusLogger())

	r.GET("/", func(c *gin.Context) {
		infoP := GetInfo()
		c.JSON(http.StatusOK, *infoP)
	})

	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	go func() {
		for {
			n := random.GetRandom(1000)
			log.WithFields(log.Fields{
				"n": n,
			}).Debug("generated random number")
			randomNumber.With(prometheus.Labels{"key": "value"}).Set(float64(n))
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		for {
			n := random.GetRandom(50)
			log.WithFields(log.Fields{
				"n": n,
			}).Debug("generated fake storageOperationMetric")
			storageOperationMetric.With(prometheus.Labels{"volume_plugin": "hongkailiu.tk/aws-ebs", "operation_name": "volume_provision"}).Observe(float64(n))
			time.Sleep(100 * time.Second)
		}
	}()

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
}
