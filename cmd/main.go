package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/leaderelection"
	rl "k8s.io/client-go/tools/leaderelection/resourcelock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/prow/pkg/interrupts"
)

var (
	leader          atomic.Value
	identity        = fmt.Sprintf("%d", os.Getpid())
	registryService RegistryService
)

type options struct {
	port                 int
	lockNamespace        string
	lockName             string
	registryLoadInterval time.Duration
}

func main() {
	opts := options{}
	flag.IntVar(&opts.port, "port", 8080, "port to serve on")
	flag.StringVar(&opts.lockNamespace, "lock-namespace", "hongkliu-test", "namespace of lock")
	flag.StringVar(&opts.lockName, "lock-name", "graph-builder-lock", "namespace of lock")
	flag.DurationVar(&opts.registryLoadInterval, "registry-load-interval", 30*time.Second, "Timeout for the operation")

	registryService = &simpleRegistryService{interval: opts.registryLoadInterval}
	registryService.Start()

	// pod name first
	podName := os.Getenv("HOSTNAME")
	if podName != "" {
		identity = podName
	}

	// Get the active kubernetes context
	cfg, err := ctrl.GetConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Error getting config")
	}

	// Create a new lock. This will be used to create a Lease resource in the cluster.
	l, err := rl.NewFromKubeconfig(
		rl.LeasesResourceLock,
		opts.lockNamespace,
		opts.lockName,
		rl.ResourceLockConfig{
			Identity: identity,
		},
		cfg,
		time.Second*10,
	)
	if err != nil {
		logrus.WithError(err).Fatal("Error creating lock")
	}

	// Create a new leader election configuration with a 15 second lease duration.
	// Visit https://pkg.go.dev/k8s.io/client-go/tools/leaderelection#LeaderElectionConfig
	// for more information on the LeaderElectionConfig struct fields
	el, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock:          l,
		LeaseDuration: time.Second * 15,
		RenewDeadline: time.Second * 10,
		RetryPeriod:   time.Second * 2,
		Name:          opts.lockName,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				logrus.Info("I am the leader!")
				if err := registryService.Load(); err != nil {
					logrus.WithError(err).Error("Error loading registry")
				}
			},
			OnStoppedLeading: func() {
				logrus.Info("I am not the leader anymore!")
			},
			OnNewLeader: func(identity string) {
				fmt.Printf("the leader is %s\n", identity)
				logrus.WithField("identity", identity).Info("The new leader is elected")
				leader.Store(identity)
			},
		},
	})
	if err != nil {
		logrus.WithError(err).Fatal("Error creating leader election")
	}

	ctx := interrupts.Context()

	// Begin the leader election process. This will block.
	go func() {
		el.Run(ctx)
	}()

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(opts.port),
		Handler: getRouter(ctx, registryService),
	}
	interrupts.ListenAndServe(server, time.Second*10)

	interrupts.WaitForGracefulShutdown()
}

func getRouter(ctx context.Context, registryService RegistryService) *http.ServeMux {
	handler := http.NewServeMux()

	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "ok")
	})

	handler.HandleFunc("/registry", func(w http.ResponseWriter, r *http.Request) {
		registry, err := registryService.Get(ctx)
		if err != nil {
			logrus.WithError(err).Error("Error loading registry data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(registry); err != nil {
			logrus.WithError(err).WithField("registry", registry).Error("failed to encode page")
		}
	})

	return handler
}

type Registry struct {
	Data string `json:"data"`
}

type RegistryService interface {
	Start()
	Get(ctx context.Context) (*Registry, error)
	Load() error
}

type simpleRegistryService struct {
	lock         sync.Mutex
	registry     *Registry
	lastModified time.Time
	interval     time.Duration
}

func (s *simpleRegistryService) Start() {
	interrupts.TickLiteral(func() {
		if err := s.Load(); err != nil {
			logrus.WithError(err).Error("Error loading registry")
		}
	}, s.interval)
}

func (s *simpleRegistryService) Load() error {
	if s.registry != nil && time.Since(s.lastModified) < s.interval {
		return nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if isLeader() {
		logrus.Info("Leader loading ...")
		// simulate the hard work
		time.Sleep(30 * time.Second)
		s.registry = &Registry{Data: "Cool registry Data"}
		s.lastModified = time.Now()
		logrus.Info("Leader loaded")
	} else {
		// TODO cache is from the leader
		return fmt.Errorf("not a leader")
	}
	return nil
}

func (s *simpleRegistryService) Get(ctx context.Context) (*Registry, error) {
	err := wait.PollUntilContextCancel(ctx, 1*time.Second, true, func(_ context.Context) (done bool, err error) {
		if s.registry == nil {
			logrus.Info("No registry available, waiting...")
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return s.registry, nil
}

func isLeader() bool {
	s, ok := leader.Load().(string)
	return ok && s == identity
}
