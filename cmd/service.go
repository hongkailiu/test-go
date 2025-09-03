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
	client       *resty.Client
	port         int
	lock         sync.Mutex
	graph        *Graph
	lastModified time.Time
	interval     time.Duration
}

func (s *simpleGraphService) Start(ctx context.Context) {
	interrupts.TickLiteral(func() {
		if err := s.Load(ctx); err != nil {
			logrus.WithError(err).WithField("interval", s.interval).Error("Error loading registry and will retry")
		}
	}, s.interval)
}

func (s *simpleGraphService) Load(_ context.Context) error {
	if s.graph != nil && time.Since(s.lastModified) < s.interval {
		return nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if isLeader() {
		logrus.Info("Leader loading ...")
		// simulate the hard work
		time.Sleep(30 * time.Second)
		s.refresh(&Graph{Data: "Cool registry Data"})
		logrus.Info("Leader loaded")
	} else {
		logrus.Info("Leader loading as non-leader ...")
		if !s.inCluster {
			return errors.New("cannot load from non-leader if not running in cluster")
		}
		leaderHostname, err := GetLeader()
		if err != nil {
			return fmt.Errorf("failed to get leader: %w", err)
		}
		registry, err := getGraphFromLeader(s.client, leaderHostname, s.port)
		if err != nil {
			return fmt.Errorf("failed to get registry from leader: %w", err)
		}
		s.refresh(registry)
		logrus.Info("Leader loaded as non-leader")
	}
	return nil
}

func getGraphFromLeader(client *resty.Client, hostname string, port int) (*Graph, error) {
	var g Graph
	res, err := client.R().
		SetResult(&g).
		Get(fmt.Sprintf("http://%s:%d/registry", hostname, port))
	if err != nil {
		return nil, fmt.Errorf("failed to get registry from leader: %w", err)
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
