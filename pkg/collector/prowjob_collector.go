package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type ProwJob struct {
	Name      string
	Namespace string
	Labels    map[string]string
}

type ProwJobController interface {
	List() []ProwJob
}

//https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
type ProwJobCollector struct {
	Name       string
	Controller ProwJobController
	log        *logrus.Logger
}

func (pjc ProwJobCollector) Describe(ch chan<- *prometheus.Desc) {
	//prometheus.DescribeByCollect(pjc, ch)
}

func (pjc ProwJobCollector) Collect(ch chan<- prometheus.Metric) {
	pjc.log.Info("ProwJobCollector collecting ,,,")
	for _, pj := range pjc.Controller.List() {
		labels := []string{"prow_job_name", "prow_job_namespace"}
		labelValues := []string{pj.Name, pj.Namespace}
		for k, v := range pj.Labels {
			labels = append(labels, k)
			labelValues = append(labelValues, v)
		}
		desc := prometheus.NewDesc(
			"prow_job_labels",
			"The number of prow jobs with the labels.",
			labels, nil,
		)
		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			float64(1),
			labelValues...,
		)
	}
}

func NewProwJobCollector(name string, prowJobController ProwJobController, reg prometheus.Registerer, log *logrus.Logger) *ProwJobCollector {
	log.Info("getting a new prow job controller ...")
	c := &ProwJobCollector{
		Name:       name,
		Controller: prowJobController,
		log:        log,
	}
	prometheus.WrapRegistererWith(prometheus.Labels{"collector_name": name}, reg).MustRegister(c)
	return c
}

func NewProwJobControllerForTest(log *logrus.Logger) *ProwJobControllerForTest {
	return &ProwJobControllerForTest{log: log}
}

type ProwJobControllerForTest struct {
	counter int
	log     *logrus.Logger
}

func (c *ProwJobControllerForTest) List() []ProwJob {
	c.log.Info("ProwJobControllerForTest list .........")
	c.counter++
	c.log.WithField("c.counter", c.counter).Info("ProwJobControllerForTest counter increased")
	if c.counter%2 == 0 {
		return []ProwJob{
			{
				Namespace: "ns2",
				Name:      "job2",
				Labels: map[string]string{
					"xyz":  "222",
					"user": "ppp",
				},
			},
		}
	}
	return []ProwJob{
		{
			Namespace: "ns1",
			Name:      "job1",
			Labels: map[string]string{
				"abc": "111",
			},
		},
		{
			Namespace: "ns3",
			Name:      "job3",
			Labels: map[string]string{
				"hello": "333",
			},
		},
	}
}
