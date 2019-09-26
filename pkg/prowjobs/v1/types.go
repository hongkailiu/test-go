package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProwJob struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

type ProwJobList struct {
	Items []ProwJob `json:"items"`
}
