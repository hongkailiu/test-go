package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"resty.dev/v3"
	"sigs.k8s.io/prow/pkg/interrupts"
)

type simpleGraphService struct {
	inCluster    bool
	podNamespace string
	client       *resty.Client
	port         int
	lock         sync.Mutex
	graph        *Graph
	lastModified time.Time
	interval     time.Duration
	myIdentity   string
}

func (s *simpleGraphService) Start(ctx context.Context) {
	interval := 30 * time.Second
	if isLeader() {
		interval = s.interval
	}

	interrupts.TickLiteral(func() {
		if err := s.Load(ctx); err != nil {
			logrus.WithError(err).WithField("myIdentity", s.myIdentity).WithField("interval", interval).Error("Error loading graph and will retry")
		}
	}, interval)
}

func (s *simpleGraphService) Load(_ context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.graph != nil && time.Since(s.lastModified) < s.interval {
		return nil
	}

	if isLeader() {
		logrus.Info("Leader loading ...")
		// simulate the hard work
		time.Sleep(30 * time.Second)
		s.refresh(&Graph{Data: "Cool graph Data"})
		logrus.Info("Leader loaded")
	} else {
		logrus.WithField("myIdentity", s.myIdentity).Info("Leader loading as non-leader ...")
		if !s.inCluster {
			return errors.New("cannot load from non-leader if not running in cluster")
		}
		leaderIdentity, err := GetLeader()
		if err != nil {
			return fmt.Errorf("failed to get leader: %w", err)
		}
		graph, err := getGraphFromLeader(s.client, s.podNamespace, leaderIdentity, s.port)
		if err != nil {
			return fmt.Errorf("failed to get graph from leader %s: %w", leaderIdentity, err)
		}
		s.refresh(graph)
		logrus.WithField("myIdentity", s.myIdentity).
			WithField("leader", leaderIdentity).
			Info("Leader loaded as non-leader")
	}
	return nil
}

func getGraphFromLeader(client *resty.Client, podNamespace, ip string, port int) (*Graph, error) {
	var g Graph
	res, err := client.R().
		SetResult(&g).
		Get(fmt.Sprintf("http://%s.%s.pod.cluster.local:%d/graph", ip, podNamespace, port))
	if err != nil {
		return nil, fmt.Errorf("failed to get graph from leader: %w", err)
	}
	if statusCode := res.StatusCode(); statusCode != 200 {
		return nil, fmt.Errorf("received an unexpected status code from from leader: %d", statusCode)
	}
	return &g, nil
}

func (s *simpleGraphService) refresh(graph *Graph) {
	s.graph = graph
	s.lastModified = time.Now()
}

func (s *simpleGraphService) Get(_ context.Context) (*Graph, error) {
	if s.graph == nil {
		return nil, fmt.Errorf("graph is still under construction")
	}
	return s.graph, nil
}
