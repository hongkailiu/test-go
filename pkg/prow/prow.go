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
	scanner := bufio.NewScanner(readCloser)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, warnKeyword) {
			ret += aurora.Sprintf("%s\n", aurora.Yellow(text))
		} else if strings.Contains(text, errorKeyword) || strings.Contains(text, fatalKeyword) {
			ret += aurora.Sprintf("%s\n", aurora.Red(text))
		}
	}
	return strings.TrimSuffix(ret, "\n"), nil
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
	c := h.clinetset
	for {
		content := content{name: name}
		h.Lock()
		h.contents[name] = &content
		h.Unlock()
		d, err := c.AppsV1().Deployments("ci").Get(name, metav1.GetOptions{})
		if err != nil {
			content.logs = aurora.Sprintf(aurora.Red(err.Error()))
			continue
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
			content.logs = err.Error()
			continue
		}
		var logs string
		for _, pod := range pods.Items {
			container := ""
			if name == "deck-internal" {
				container = "deck"
			}
			if name == "boskos" {
				container = "boskos"
			}
			lines, err := containerLog(h.clinetset, "ci", pod.Name, container)
			if err != nil {
				lines = aurora.Sprintf(aurora.Red(err.Error()))
			}
			logs += lines
		}
		content.logs = logs
		time.Sleep(60 * time.Second)
	}
}

func (h *memoryOutputHandler) display() error {
	writer := uilive.New()
	writer.Start()
	for {
		var lines string
		h.RLock()
		keys := make([]string, 0, len(h.contents))
		for k := range h.contents {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			content := h.contents[k]
			lines += fmt.Sprintf("%s\n", content.Header())
			if content.logs != "" {
				lines += fmt.Sprintf("%s\n", content.logs)
			}
		}
		h.RUnlock()
		if _, err := fmt.Fprintf(writer, lines); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
}
