package http

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed.",
		},
		[]string{"path", "hostname", "method"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "http request duration seconds",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"path", "hostname", "method", "status"},
	)

	randomNumber = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "random_number",
			Help: "the value of random number.",
		}, []string{"key", "hostname"},
	)

	//https://github.com/kubernetes/kubernetes/blob/master/pkg/volume/util/metrics.go
	storageOperationMetric = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "storage_operation_duration_seconds",
			Help:    "Storage operation duration",
			Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10, 15, 25, 50},
		},
		[]string{"volume_plugin", "operation_name", "hostname"},
	)

	metricLabels = []string{
		// namespace of the job
		"job_namespace",
		// name of the job
		"job_name",
		// type of the prowjob: presubmit, postsubmit, periodic, batch
		"type",
		// state of the prowjob: triggered, pending, success, failure, aborted, error
		"state",
		// the org of the prowjob's repo
		"org",
		// the prowjob's repo
		"repo",
		// the base_ref of the prowjob's repo
		"base_ref",
	}

	prowJobTransitions = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "prowjob_state_transitions",
		Help: "Number of prowjobs transitioning states",
	}, metricLabels)
)

func mustRegister(component, selector string, client prowJobClient) *prometheus.Registry {
	registry := prometheus.NewRegistry()
	prometheus.WrapRegistererWith(prometheus.Labels{"collector_name": component}, registry).MustRegister(&prowJobCollector{
		name:     component,
		client:   client,
		selector: selector,
	})
	registry.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	registry.MustRegister(httpRequestsTotal)
	registry.MustRegister(httpRequestDuration)
	registry.MustRegister(randomNumber)
	registry.MustRegister(storageOperationMetric)
	registry.MustRegister(prowJobTransitions)
	return registry
}

type jobLabel struct {
	jobNamespace string
	jobName      string
	jobType      string
	state        string
	org          string
	repo         string
	baseRef      string
}

func (jl *jobLabel) values() []string {
	return []string{jl.jobNamespace, jl.jobName, jl.jobType, jl.state, jl.org, jl.repo, jl.baseRef}
}

func generateJobLabel(n int) jobLabel {
	switch n % 3 {
	case 0:
		return jobLabel{
			jobNamespace: "ns1",
			jobName:      "job1",
			jobType:      "presubmit",
			state:        "failure",
			org:          "codeready-toolchain",
			repo:         "host-operator",
			baseRef:      "master",
		}
	case 1:
		return jobLabel{
			jobNamespace: "ns2",
			jobName:      "job2",
			jobType:      "periodic",
			state:        "success",
			org:          "operator-framework",
			repo:         "operator-sdk",
			baseRef:      "master",
		}
	case 2:
		return jobLabel{
			jobNamespace: "ns3",
			jobName:      "job3",
			jobType:      "postsubmit",
			state:        "failure",
			org:          "openshift",
			repo:         "openshift-azure",
			baseRef:      "master",
		}
	}
	return jobLabel{}
}
