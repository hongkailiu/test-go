package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProwJobAgent specifies the controller (such as plank or jenkins-agent) that runs the job.
type ProwJobAgent string

const (
	// KubernetesAgent means prow will create a pod to run this job.
	KubernetesAgent ProwJobAgent = "kubernetes"
	// JenkinsAgent means prow will schedule the job on jenkins.
	JenkinsAgent ProwJobAgent = "jenkins"
	// KnativeBuildAgent means prow will schedule the job via a build-crd resource.
	KnativeBuildAgent ProwJobAgent = "knative-build"
	// TektonAgent means prow will schedule the job via a tekton PipelineRun CRD resource.
	TektonAgent = "tekton-pipeline"
)

type ProwJob struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProwJobSpec `json:"spec,omitempty"`
}

type ProwJobSpec struct {
	Agent ProwJobAgent `json:"agent,omitempty"`
}

type ProwJobList struct {
	Items []ProwJob `json:"items"`
}
