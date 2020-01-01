package prow

import (
	"bufio"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gosuri/uilive"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"

	"github.com/hongkailiu/test-go/pkg/ocutil"
	cmdconfig "github.com/hongkailiu/test-go/pkg/testctl/cmd/config"
)

func Monitor(pc *cmdconfig.ProwConfig) error {
	logrus.Debugf("Start monitoring Prow deployments")
	kcPath := pc.KubeConfigPath
	if kcPath == "" {
		kcPath = filepath.Join(homedir.HomeDir(), ".kube", "config")
		logrus.Debugf("--kubeconfig is unset, using '%s'", kcPath)
	}
	cfgs, currentContext, err := ocutil.LoadKubeConfigs(kcPath)
	logrus.Debugf("The current context is '%s'", currentContext)
	if err != nil {
		return err
	}
	config := cfgs[currentContext]
	clientset, err := kubernetes.NewForConfig(&config)
	if err != nil {
		return err
	}
	deployments, err := clientset.AppsV1().Deployments("ci").List(metav1.ListOptions{LabelSelector: "app=prow"})
	if err != nil {
		return err
	}
	var handler outputHandler
	handler = newMemoryOutputHandler(clientset)
	for _, d := range deployments.Items {
		logrus.WithField("d.Name", d.Name).Debugf("found d")
		go handler.getAndSave(d.Name)
	}
	if err := handler.display(); err != nil {
		return err
	}
	return nil
}

type content struct {
	name      string
	version   string
	current   int32
	desired   int32
	updated   int32
	available int32
	logs      string
}

func (c *content) Header() string {
	message := fmt.Sprintf("%s at %s [%d/%d]", c.name, c.version, c.current, c.desired)
	if c.updated != c.desired {
		message += fmt.Sprintf(" (%d stale replicas)", c.desired-c.updated)
	}
	if c.available != c.desired {
		message += fmt.Sprintf(" (%d unavailable replicas)", c.desired-c.available)
	}
	message += ":"
	message = aurora.Sprintf(aurora.Bold(message))
	if c.desired != c.current {
		message = aurora.Sprintf(aurora.Red(message))
	}
	if len(c.logs) == 0 {
		message = aurora.Sprintf("%s %s", message, aurora.Green("OK"))
	}
	return message
}

const (
	warnKeyword  = `"level":"warning"`
	errorKeyword = `"level":"error"`
	fatalKeyword = `"level":"fatal"`
)

func renderFlavor(clientset *kubernetes.Clientset, project, podName, dc string) (string, error) {
	var lines []string
	pod, err := clientset.CoreV1().Pods(project).Get(podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if pod.Status.Phase != corev1.PodRunning {
		text := fmt.Sprintf("pod %s is %s: %s, %s", pod.Name, pod.Status.Phase, pod.Status.Reason, pod.Status.Message)
		line := aurora.Sprintf(aurora.Yellow(text))
		if sets.NewString("Failed", "Unknown", "CrashLoopBackOff").Has(string(pod.Status.Phase)) {
			line = aurora.Sprintf(aurora.Red(text))
		}
		lines = append(lines, line)
	}

	for _, container := range pod.Status.ContainerStatuses {
		if container.Name == dc {
			if container.State.Running == nil {
				if container.State.Waiting != nil {
					lines = append(lines, aurora.Sprintf(aurora.Yellow(fmt.Sprintf("pod %s is waitng: %s", podName, container.State.Waiting.Reason))))
					lines = append(lines, aurora.Sprintf(aurora.Yellow(fmt.Sprintf("\t%s", container.State.Waiting.Message))))
				}
				if container.State.Terminated != nil {
					lines = append(lines, aurora.Sprintf(aurora.Red(fmt.Sprintf("pod %s is waitng: %s", podName, container.State.Terminated.Reason))))
					lines = append(lines, aurora.Sprintf(aurora.Red(fmt.Sprintf("\t%s", container.State.Terminated.Message))))
				}
			}
		}
		if container.RestartCount != 0 {
			lines = append(lines, aurora.Sprintf(aurora.Red(fmt.Sprintf("pod %s has restarted %d times", podName, container.RestartCount))))
		}
	}
	return strings.Join(lines, "\n"), nil
}

func containerLog(clientset *kubernetes.Clientset, projectName string, podName string, container string) (ret string, returnedErr error) {
	since := int64(20 * 60)
	logOptions := &corev1.PodLogOptions{
		Container:    container,
		SinceSeconds: &since,
	}
	readCloser, err := clientset.CoreV1().Pods(projectName).GetLogs(podName, logOptions).Stream()
	if err != nil {
		return "", err
	}
	defer func() {
		err := readCloser.Close()
		if err != nil {
			returnedErr = err
		}
	}()
	var lines []string
	scanner := bufio.NewScanner(readCloser)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, warnKeyword) {
			lines = append(lines, aurora.Sprintf(aurora.Yellow(text)))
		} else if strings.Contains(text, errorKeyword) || strings.Contains(text, fatalKeyword) {
			lines = append(lines, aurora.Sprintf(aurora.Red(text)))
		}
	}
	if len(lines) > 5 {
		lines = lines[len(lines)-5:]
	}
	return strings.Join(lines, "\n"), nil
}

type outputHandler interface {
	getAndSave(name string)
	display() error
}

type memoryOutputHandler struct {
	clinetset *kubernetes.Clientset
	contents  map[string]*content
	sync.RWMutex
}

func newMemoryOutputHandler(clientset *kubernetes.Clientset) *memoryOutputHandler {
	ret := memoryOutputHandler{
		clinetset: clientset,
		contents:  map[string]*content{},
	}
	return &ret
}

func (h *memoryOutputHandler) getAndSave(name string) {
	for {
		content := h.getContent(name)
		h.Lock()
		h.contents[name] = content
		h.Unlock()
		time.Sleep(60 * time.Second)
	}
}

func (h *memoryOutputHandler) getContent(name string) *content {
	c := h.clinetset
	content := &content{name: name}
	d, err := c.AppsV1().Deployments("ci").Get(name, metav1.GetOptions{})
	if err != nil {
		content.logs = aurora.Sprintf("Failed to get deployment '%s': '%s'", name, aurora.Red(err.Error()))
		return content
	}
	content.desired = *d.Spec.Replicas
	content.current = d.Status.Replicas
	content.updated = d.Status.UpdatedReplicas
	content.available = d.Status.AvailableReplicas
	content.version = "<unknown-version>"
	for _, container := range d.Spec.Template.Spec.Containers {
		containerName := name
		if name == "boskos-metrics" {
			containerName = "metrics"
		}
		if name == "jenkins-dev-operator" {
			containerName = "jenkins-operator"
		}
		if name == "deck-internal" {
			containerName = "deck"
		}
		if container.Name == containerName {
			parts := strings.Split(container.Image, ":")
			content.version = parts[len(parts)-1]
		}
	}
	pods, err := c.CoreV1().Pods("ci").List(metav1.ListOptions{LabelSelector: fmt.Sprintf("component=%s", name)})
	if err != nil {
		content.logs = aurora.Sprintf("Failed to list pod of component '%s': '%s'", name, aurora.Red(err.Error()))
		return content
	}
	var logs []string
	for _, pod := range pods.Items {
		lines, err := renderFlavor(h.clinetset, "ci", pod.Name, name)
		if err != nil {
			lines = aurora.Sprintf("Failed to render flavor: '%s'", aurora.Red(err.Error()))
		}
		if lines != "" {
			logs = append(logs, lines)
		}
		container := ""
		if name == "deck-internal" {
			container = "deck"
		}
		if name == "boskos" {
			container = "boskos"
		}
		lines, err = containerLog(h.clinetset, "ci", pod.Name, container)
		if err != nil {
			lines = aurora.Sprintf("Failed to get container log: '%s'", aurora.Red(err.Error()))
		}
		if lines != "" {
			logs = append(logs, lines)
		}
	}
	content.logs = strings.Join(logs, "\n")
	return content
}

func (h *memoryOutputHandler) display() error {
	writer := uilive.New()
	writer.Start()
	for {
		if _, err := fmt.Fprintf(writer, fmt.Sprintf("%s\n", h.displayString())); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
}

func (h *memoryOutputHandler) displayString() string {
	var lines []string
	//h.RLock()
	//defer h.RUnlock()
	keys := make([]string, 0, len(h.contents))
	for k := range h.contents {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		content := h.contents[k]
		lines = append(lines, content.Header())
		if content.logs != "" {
			lines = append(lines, content.logs)
		}
	}
	return strings.Join(lines, "\n")
}
