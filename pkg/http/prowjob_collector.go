/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package http

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	prowapi "github.com/hongkailiu/test-go/pkg/prowjobs/v1"
)

type prowJobClient interface {
	List(opts metav1.ListOptions) (*prowapi.ProwJobList, error)
}

//https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
type prowJobCollector struct {
	name     string
	client   prowJobClient
	selector string
}

func (pjc prowJobCollector) Describe(ch chan<- *prometheus.Desc) {
	//prometheus.DescribeByCollect(pjc, ch) //Normally, we would need this line.
	//https://godoc.org/github.com/prometheus/client_golang/prometheus#hdr-Custom_Collectors_and_constant_Metrics
	//This is a take-our-own-risk action and also a compromise for implementing a metric with
	//both dynamic keys and dynamic values in the label set.
	//The formed scraping data have to be tested locally with prometheus server.
	//Another option is taken for kube_pod_labels: https://github.com/kubernetes/kube-state-metrics/blob/master/docs/pod-metrics.md
	//If we follow the implementation there, we would need to import pkgs from "k8s.io/kube-state-metrics".
}

func (pjc prowJobCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debug("ProwJobCollector collecting ...")
	prowJobs, err := pjc.client.List(metav1.ListOptions{LabelSelector: pjc.selector})
	if err != nil {
		logrus.WithError(err).Errorf("Failed to list prow jobs with selector '%s'", pjc.selector)
		return
	}
	latestJobs := getLatest(prowJobs.Items)
	for _, pj := range latestJobs {
		agent := string(pj.Spec.Agent)
		pjLabelKeys, pjLabelValues := kubeLabelsToPrometheusLabels(pj.Labels, "label_")
		pjLabelKeys = append([]string{"prow_job_name", "prow_job_namespace", "prow_job_agent"}, pjLabelKeys...)
		pjLabelValues = append([]string{pj.Spec.Job, pj.Namespace, agent}, pjLabelValues...)
		labelDesc := prometheus.NewDesc(
			"prow_job_labels",
			"Kubernetes labels converted to Prometheus labels.",
			pjLabelKeys, nil,
		)
		ch <- prometheus.MustNewConstMetric(
			labelDesc,
			prometheus.GaugeValue,
			// always 1 since there is only 1 prow job for each namespace and each prow job name
			float64(1),
			pjLabelValues...,
		)
		pjAnnotationKeys, pjAnnotationValues := kubeLabelsToPrometheusLabels(pj.Annotations, "annotation_")
		pjAnnotationKeys = append([]string{"prow_job_name", "prow_job_namespace", "prow_job_agent"}, pjAnnotationKeys...)
		pjAnnotationValues = append([]string{pj.Name, pj.Namespace, agent}, pjAnnotationValues...)
		annotationDesc := prometheus.NewDesc(
			"prow_job_annotations",
			"Kubernetes annotations converted to Prometheus labels.",
			pjAnnotationKeys, nil,
		)
		ch <- prometheus.MustNewConstMetric(
			annotationDesc,
			prometheus.GaugeValue,
			float64(1),
			pjAnnotationValues...,
		)
	}
}

var (
	invalidLabelCharRE    = regexp.MustCompile(`[^a-zA-Z0-9_]`)
	escapeWithDoubleQuote = strings.NewReplacer("\\", `\\`, "\n", `\n`, "\"", `\"`)
)

// https://github.com/kubernetes/kube-state-metrics/blob/1d69c1e637564aec4591b5b03522fa8b5fca6597/internal/store/utils.go#L60
func kubeLabelsToPrometheusLabels(labels map[string]string, prefix string) ([]string, []string) {
	labelKeys := make([]string, 0, len(labels))
	for k := range labels {
		labelKeys = append(labelKeys, k)
	}
	sort.Strings(labelKeys)

	labelValues := make([]string, 0, len(labels))
	for i, k := range labelKeys {
		labelKeys[i] = fmt.Sprintf("%s%s", prefix, sanitizeLabelName(k))
		labelValues = append(labelValues, escapeString(labels[k]))
	}
	return labelKeys, labelValues
}

func sanitizeLabelName(s string) string {
	return invalidLabelCharRE.ReplaceAllString(s, "_")
}

// https://github.com/kubernetes/kube-state-metrics/blob/1d69c1e637564aec4591b5b03522fa8b5fca6597/pkg/metric/metric.go#L96
func escapeString(v string) string {
	return escapeWithDoubleQuote.Replace(v)
}

func getLatest(jobs []prowapi.ProwJob) map[string]prowapi.ProwJob {
	latest := map[string]time.Time{}
	latestJobs := map[string]prowapi.ProwJob{}
	for _, job := range jobs {
		if _, ok := latest[job.Spec.Job]; !ok {
			latest[job.Spec.Job] = job.Status.StartTime.Time
			latestJobs[job.Spec.Job] = job
			continue
		}
		if job.Status.StartTime.Time.After(latest[job.Spec.Job]) {
			latest[job.Spec.Job] = job.Status.StartTime.Time
			latestJobs[job.Spec.Job] = job
		}
	}
	return latestJobs
}
