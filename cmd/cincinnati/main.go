package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"sigs.k8s.io/prow/pkg/interrupts"

	"github.com/hongkailiu/test-go/pkg/cincinnati"
)

// TODO: accept config file

type options struct {
	Address        string
	MetricsAddress string
	Registry       string
	Repo           string
	GraphDataDir   string

	MockDir        string
	DataDir        string
	GraphFile      string
	GracePeriod    time.Duration
	MaxConcurrency int
	LogLevel       string
}

var opts options

var ctx = interrupts.Context()

var rootCmd = &cobra.Command{
	Use:   "cincinnati",
	Short: "A Cincinnati update graph server",
	Long:  "cincinnati is a GoLang implementation of the Red Hat OpenShift Cincinnati update graph protocol",
	Run: func(cmd *cobra.Command, args []string) {
		level, err := logrus.ParseLevel(opts.LogLevel)
		if err != nil {
			logrus.WithError(err).WithField("level", opts.LogLevel).Fatal("invalid log level")
		}
		logrus.SetLevel(level)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		logrus.SetReportCaller(true)

		client := retryablehttp.NewClient()
		client.HTTPClient.Timeout = 30 * time.Second
		client.RetryMax = 3
		client.RetryWaitMin = 1 * time.Second
		client.RetryWaitMax = 5 * time.Second

		l := logrus.New()
		l.SetLevel(logrus.WarnLevel)
		l.WithField("subComponent", "retryablehttp")
		client.Logger = l

		repo := cincinnati.NewRepo(client.StandardClient(),
			opts.Registry, opts.Repo, opts.DataDir, opts.MockDir, opts.MaxConcurrency)

		c := cache.New(5*time.Minute, 10*time.Minute)

		gb := cincinnati.NewGraphBuilder(opts.GraphFile, opts.GraphDataDir, opts.MockDir, c, repo)
		if err := gb.Start(ctx); err != nil {
			logrus.WithError(err).Fatal("Failed to start server")
		}

		r := gin.Default()
		p := ginprometheus.NewWithConfig(ginprometheus.Config{
			Subsystem:          cincinnati.MetricsPrefix,
			DisableBodyReading: true,
		})
		r.Use(p.HandlerFunc())

		server := &http.Server{
			Addr:    opts.Address,
			Handler: cincinnati.GetHandler(r, gb),
		}

		go func() {
			if err := startServer(server); err != nil {
				logrus.WithError(err).Fatal("Failed to start server")
			}
		}()

		metricsRouter := gin.New()
		metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))

		metricsServer := &http.Server{
			Addr:    opts.MetricsAddress,
			Handler: metricsRouter.Handler(),
		}

		go func() {
			if err := startServer(metricsServer); err != nil {
				logrus.WithError(err).Fatal("Failed to start metrics server")
			}
		}()

		waitForShutdown(opts.GracePeriod, server, metricsServer)
		logrus.Info("Process language gracefully")
	},
}

func waitForShutdown(period time.Duration, servers ...*http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logrus.Info("Shutdown signal received")

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), period)
	defer cancel()

	for i, server := range servers {
		if err := server.Shutdown(ctx); err != nil {
			logrus.WithError(err).WithField("i", i).Error("Server forced to shutdown")
		}
		logrus.WithField("i", i).Info("Server shutdown gracefully")
	}
}

func startServer(server *http.Server) error {
	logrus.Info("Start server")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func init() {
	rootCmd.Flags().StringVar(&opts.Address, "address", cincinnati.DefaultPort, "Address to run the server with")
	rootCmd.Flags().StringVar(&opts.MockDir, "mock-dir", "", "Path to the directory containing mock files")
	rootCmd.Flags().StringVar(&opts.DataDir, "data-dir", "data", "Path to the directory containing image info files")
	rootCmd.Flags().StringVar(&opts.GraphDataDir, "graph-data-dir", "/tmp/cincinnati/graph-data",
		"Path to the directory containing graph data")
	rootCmd.Flags().StringVar(&opts.Registry, "registry", "https://quay.io", "Registry URL")
	rootCmd.Flags().StringVar(&opts.Repo, "repo", "openshift-release-dev/ocp-release", "Repo in form of org/repo")
	rootCmd.Flags().StringVar(&opts.GraphFile, "graph-file", "data/graph.json.gz", "Graph file path")
	rootCmd.Flags().DurationVar(&opts.GracePeriod, "gracePeriod", time.Second*10, "Grace period for server shutdown")
	rootCmd.Flags().IntVar(&opts.MaxConcurrency, "max-concurrency", cincinnati.DefaultMaxConcurrency,
		"Maximum number of concurrent in-flight goroutines to scrape the registry")
	rootCmd.Flags().StringVar(&opts.LogLevel, "log-level", "info", "Set log level (debug, info, warn, error)")
	rootCmd.Flags().StringVar(&opts.MetricsAddress, "metrics-address", cincinnati.DefaultMetricsPort, "Address to run the metrics server with")

	if v := os.Getenv("CINCINNATI_REGISTRY"); v != "" {
		opts.Registry = v
	}

	if v := os.Getenv("CINCINNATI_REPO"); v != "" {
		opts.Repo = v
	}

	if v := os.Getenv("MOCK_DIR"); v != "" {
		opts.MockDir = v
	}
}

func main() {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		os.Exit(1)
	}
}
