package ocpsanity

import (
	"fmt"

	"github.com/hongkailiu/test-go/pkg/ocutil"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// VERSION of the ocptf cmd
	VERSION = "0.0.8"
)

var (
	oc *ocutil.CLI
)

type SanitySummary struct {
	ProjectTotal       int              `yaml:"projectTotal"`
	ProjectSummaryList []ProjectSummary `yaml:"projectSummaryList"`
}

type ProjectSummary struct {
	ProjectName               string `yaml:"projectName"`
	DCTotal                   int    `yaml:"dcTotal"`
	DCReplicaNotSatisfied     int    `yaml:"dcReplicaNotSatisfied"`
	DeployTotal               int    `yaml:"deployTotal"`
	DeployReplicaNotSatisfied int    `yaml:"deployReplicaNotSatisfied"`
	STSTotal                  int    `yaml:"stsTotal"`
	STSReplicaNotSatisfied    int    `yaml:"stsReplicaNotSatisfied"`
	DSTotal                   int    `yaml:"dsTotal"`
	DSDesireNotSatisfied      int    `yaml:"dsDesireNotSatisfied"`

	PodTotal          int `yaml:"podTotal"`
	PodNotRunning     int `yaml:"podNotRunning"`
	ContainerNotReady int `yaml:"containerNotReady"`
}

func StartSanityCheck(configPath string) error {
	oc = ocutil.NewCLI(configPath)

	sanitySummary := SanitySummary{ProjectTotal: 0, ProjectSummaryList: []ProjectSummary{}}

	projectList, err := oc.ProjectClient().Projects().List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	sanitySummary.ProjectTotal = len(projectList.Items)

	for _, project := range projectList.Items {
		projectSummary := ProjectSummary{}
		projectSummary.ProjectName = project.Name
		log.WithFields(log.Fields{"name": project.Name}).Info("Handle project")
		deployConfigList, err := oc.DeployConfigClient().DeploymentConfigs(project.Name).List(metav1.ListOptions{})
		if err != nil {
			return err
		}

		projectSummary.DCTotal = len(deployConfigList.Items)
		for _, dc := range deployConfigList.Items {
			if dc.Spec.Replicas != dc.Status.Replicas {
				log.WithFields(log.Fields{"project": project.Name, "dc": dc.Name, "dc.Spec.Replicas": dc.Spec.Replicas,
					"dc.Status.Replicas": dc.Status.Replicas}).
					Warn("Handle deploymentconfig.apps.openshift.io: Replicas not satisfied")
				projectSummary.DCReplicaNotSatisfied++
			} else {
				log.WithFields(log.Fields{"project": project.Name, "dc": dc.Name, "dc.Spec.Replicas": dc.Spec.Replicas,
					"dc.Status.Replicas": dc.Status.Replicas}).
					Info("Handle deploymentconfig.apps.openshift.io")
			}
		}

		deploymentList, err := oc.K8SClientSet().AppsV1().Deployments(project.Name).List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		projectSummary.DeployTotal = len(deploymentList.Items)
		for _, d := range deploymentList.Items {
			if *d.Spec.Replicas != d.Status.Replicas {
				log.WithFields(log.Fields{"project": project.Name, "dc": d.Name, "d.Spec.Replicas": *d.Spec.Replicas,
					"d.Status.Replicas": d.Status.Replicas}).
					Warn("Handle deployment: Replicas not satisfied")
				projectSummary.DeployReplicaNotSatisfied++
			} else {
				log.WithFields(log.Fields{"project": project.Name, "dc": d.Name, "d.Spec.Replicas": *d.Spec.Replicas,
					"d.Status.Replicas": d.Status.Replicas}).
					Info("Handle deployment")
			}
		}

		statefulSetList, err := oc.K8SClientSet().AppsV1().StatefulSets(project.Name).List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		projectSummary.STSTotal = len(statefulSetList.Items)
		for _, ss := range statefulSetList.Items {
			if *ss.Spec.Replicas != ss.Status.Replicas {
				log.WithFields(log.Fields{"project": project.Name, "ss": ss.Name, "ss.Spec.Replicas": *ss.Spec.Replicas,
					"ss.Status.Replicas": ss.Status.Replicas}).
					Warn("Handle sts: Replicas not satisfied")
				projectSummary.STSReplicaNotSatisfied++
			} else {
				log.WithFields(log.Fields{"project": project.Name, "ss": ss.Name, "ss.Spec.Replicas": *ss.Spec.Replicas,
					"ss.Status.Replicas": ss.Status.Replicas}).
					Info("Handle sts")
			}
		}

		daemonSetList, err := oc.K8SClientSet().AppsV1().DaemonSets(project.Name).List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		projectSummary.DSTotal = len(daemonSetList.Items)
		for _, ds := range daemonSetList.Items {
			if ds.Status.DesiredNumberScheduled != ds.Status.CurrentNumberScheduled {
				log.WithFields(log.Fields{"project": project.Name, "ds": ds.Name,
					"ds.Status.DesiredNumberScheduled": ds.Status.DesiredNumberScheduled,
					"ds.Status.CurrentNumberScheduled": ds.Status.CurrentNumberScheduled}).
					Warn("Handle daemon set: Not scheduled")
				projectSummary.DSDesireNotSatisfied++
			} else {
				log.WithFields(log.Fields{"project": project.Name, "ds": ds.Name,
					"ds.Status.DesiredNumberScheduled": ds.Status.DesiredNumberScheduled,
					"ds.Status.CurrentNumberScheduled": ds.Status.CurrentNumberScheduled}).
					Warn("Handle daemon set")
			}
		}

		PodList, err := oc.K8SClientSet().CoreV1().Pods(project.Name).List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		projectSummary.PodTotal = len(PodList.Items)
		for _, pod := range PodList.Items {
			if pod.Status.Phase != corev1.PodRunning {
				log.WithFields(log.Fields{"project": project.Name, "pod": pod.Name, "pod.Status.Phase": pod.Status.Phase,
					"pod.Status.Reason": pod.Status.Reason}).
					Warn("Handle pod: Not Running")
				projectSummary.PodNotRunning++
			} else {
				for _, cs := range pod.Status.ContainerStatuses {
					if !cs.Ready {
						log.WithFields(log.Fields{"project": project.Name, "pod": pod.Name,
							"pod.Status.Phase": pod.Status.Phase, "cs.Name": cs.Name,
							"cs.ContainerID": cs.ContainerID, "cs.Ready": cs.Ready}).
							Warn("Handle pod: Not Ready")
						projectSummary.ContainerNotReady++
					} else {
						log.WithFields(log.Fields{"project": project.Name, "pod": pod.Name,
							"pod.Status.Phase": pod.Status.Phase, "cs.Name": cs.Name,
							"cs.ContainerID": cs.ContainerID, "cs.Ready": cs.Ready}).
							Info("Handle pod")
					}
				}

			}
		}

		sanitySummary.ProjectSummaryList = append(sanitySummary.ProjectSummaryList, projectSummary)
	}

	bytes, err := yaml.Marshal(&sanitySummary)
	if err != nil {
		return err
	}
	fmt.Printf("--- sanitySummary --- \n%s\n", string(bytes))

	return nil
}
