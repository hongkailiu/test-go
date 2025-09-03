package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/leaderelection"
	rl "k8s.io/client-go/tools/leaderelection/resourcelock"
	"resty.dev/v3"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/prow/pkg/interrupts"
)

var (
	leader       atomic.Value
	myIdentity   = fmt.Sprintf("%d", os.Getpid())
	graphService GraphService
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
	flag.DurationVar(&opts.registryLoadInterval, "registry-load-interval", 30*time.Minute, "Timeout for the operation")

	flag.Parse()

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	// pod name first
	podName := os.Getenv("HOSTNAME")
	if podName != "" {
		myIdentity = podName
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
			Identity: myIdentity,
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
				if err := graphService.Load(ctx); err != nil {
					logrus.WithError(err).Error("Error loading registry")
				}
			},
			OnStoppedLeading: func() {
				logrus.Info("I am not the leader anymore!")
			},
			OnNewLeader: func(identity string) {
				leader.Store(identity)
				logrus.WithField("isLeader", isLeader()).
					WithField("myIdentity", myIdentity).
					WithField("identity", identity).
					Info("The new leader is elected")

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

	if err := wait.PollUntilContextTimeout(ctx, time.Second, time.Minute, true, func(_ context.Context) (done bool, err error) {
		_, err = GetLeader()
		if err != nil {
			if errors.Is(err, NoLeaderError) {
				return false, nil
			}
			return false, fmt.Errorf("failed to get leader: %w", err)
		}
		return true, nil
	}); err != nil {
		logrus.WithError(err).Fatal("Error getting leader")
	}

	client := resty.New()
	defer func() {
		if err := client.Close(); err != nil {
			logrus.WithError(err).Error("failed to close client")
		}
	}()

	graphService = &simpleGraphService{interval: opts.registryLoadInterval, port: opts.port, client: client, inCluster: podName != ""}
	graphService.Start(ctx)

	server := &http.Server{
		Addr: ":" + strconv.Itoa(opts.port),
		//Addr:    ":" + strconv.Itoa(8081),
		Handler: getRouter(ctx, graphService),
	}
	interrupts.ListenAndServe(server, time.Second*10)
	interrupts.WaitForGracefulShutdown()
}

type Graph struct {
	Data string `json:"data"`
}

type GraphService interface {
	Start(ctx context.Context)
	Get(ctx context.Context) (*Graph, error)
	Load(ctx context.Context) error
}

func isLeader() bool {
	s, ok := leader.Load().(string)
	return ok && s == myIdentity
}

var NoLeaderError = errors.New("no leader found")

func GetLeader() (string, error) {
	if s, ok := leader.Load().(string); ok {
		return s, nil
	}
	return "", NoLeaderError
}
