package http

import (
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
	if c.counter%2 == 0 {
		return &prowapi.ProwJobList{
			Items: []prowapi.ProwJob{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns2",
						Name:      "job2",
						Labels: map[string]string{
							"xyz":  "222",
							"user": "ppp",
						},
					},
				},
			},
		}, nil
	}

	return &prowapi.ProwJobList{
		Items: []prowapi.ProwJob{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns1",
					Name:      "job1",
					Labels: map[string]string{
						"abc": "111",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns3",
					Name:      "job3",
					Labels: map[string]string{
						"hello": "333",
					},
				},
			},
		},
	}, nil

}

func NewProwJobControllerForTest(log *logrus.Logger) *ProwJobControllerForTest {
	return &ProwJobControllerForTest{log: log}
}
