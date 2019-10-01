package http

import (
	"fmt"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	prowapi "github.com/hongkailiu/test-go/pkg/prowjobs/v1"
)

type ProwJobControllerForTest struct {
	counter int
	log     *logrus.Logger
}

func (c *ProwJobControllerForTest) List(opts metav1.ListOptions) (*prowapi.ProwJobList, error) {
	c.log.Info("ProwJobControllerForTest list .........")
	c.counter++
	c.log.WithField("c.counter", c.counter).Info("ProwJobControllerForTest counter increased")
	switch c.counter % 3 {
	case 0:
		return &prowapi.ProwJobList{
			Items: []prowapi.ProwJob{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns2",
						Name:      "job2",
						Labels: map[string]string{
							"created-by-prow": "true",
							"user":            "ppp",
						},
						Annotations: map[string]string{
							"prow.k8s.io/job": "branch-ci-openshift-release-master-core-apply",
						},
					},
				},
			},
		}, nil
	case 1:
		return &prowapi.ProwJobList{
			Items: []prowapi.ProwJob{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "job1",
						Labels: map[string]string{
							"event-GUID": "bcb9fc8c-e3b3-11e9-98b0-2da0529ae97c",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns3",
						Name:      "job3",
						Labels: map[string]string{
							"prow.k8s.io/refs.repo": "release",
							"prow.k8s.io/refs.org":  "openshift",
						},
						Annotations: map[string]string{
							"prow.k8s.io/job": "branch-ci-openshift-release-master-core-apply",
						},
					},
				},
			},
		}, nil
	case 2:
		return &prowapi.ProwJobList{
			Items: []prowapi.ProwJob{
			},
		}, nil

	}

	return nil, fmt.Errorf("this should have never happened")

}

func NewProwJobControllerForTest(log *logrus.Logger) *ProwJobControllerForTest {
	return &ProwJobControllerForTest{log: log}
}
