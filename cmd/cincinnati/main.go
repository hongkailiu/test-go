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
	"github.com/oklog/run"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	ginprometheus "github.com/zsais/go-gin-prometheus"

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

		var g run.Group
		ctx, cancel := context.WithCancel(context.Background())

		{
			// Do the initial work once and the work will be canceled with the context
			logger := logrus.WithField("worker", "init")
			go func(ctx context.Context) {
				logger.Info("Initial work started")
				if err := gb.CacheGraphData(); err != nil {
					logger.WithError(err).Fatal("Work completed with error")
				}
				if err := gb.CacheGraph(ctx); err != nil {
					logger.WithError(err).Fatal("Work completed with error")
				}
				logger.Println("Initial work completed")
			}(ctx)
		}

		{
			logger := logrus.WithField("worker", "CacheGraphData")
			g.Add(func() error {

				logger.Info("Worker started")
				ticker := time.NewTicker(time.Minute)
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						logger.Info("Worker received shutdown signal")
						return nil
					case t := <-ticker.C:
						logger.WithField("", t).Println("Work started")
						if err := gb.CacheGraphData(); err != nil {
							logger.WithError(err).Error("Work completed with error")
						}
						logger.Println("Work completed")
					}
				}
			}, func(err error) {
				logger.Warn("Worker stopping")
				cancel()
			})
		}

		{
			logger := logrus.WithField("worker", "CacheGraph")
			g.Add(func() error {

				logger.Info("Worker started")
				ticker := time.NewTicker(2 * time.Hour)
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						logger.Info("Worker received shutdown signal")
						return nil
					case t := <-ticker.C:
						logger.WithField("", t).Println("Work started")
						if err := gb.CacheGraph(ctx); err != nil {
							logrus.WithError(err).Error("Work completed with error")
						}
						logger.Println("Work completed")
					}
				}
			}, func(err error) {
				logger.Warn("Worker stopping")
				cancel()
			})
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

		metricsRouter := gin.New()
		metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))

		metricsServer := &http.Server{
			Addr:    opts.MetricsAddress,
			Handler: metricsRouter.Handler(),
		}

		{
			logger := logrus.WithField("server", "main")
			g.Add(func() error {
				logger.Info("Server started")
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.WithError(err).Error("Server stopped unexpectedly")
					return err
				}
				logger.Info("Server stopped")
				return nil
			}, func(err error) {
				logger.Info("Server stopping")
				shutdownCtx, cancelFn := context.WithTimeout(context.Background(), opts.GracePeriod)
				defer cancelFn()
				if err := server.Shutdown(shutdownCtx); err != nil {
					logger.WithError(err).Error("Server shutdown failed")
					return
				}
				logger.Info("Server stopped gracefully")
			})
		}

		{
			logger := logrus.WithField("server", "metrics")
			g.Add(func() error {
				logger.Info("Server started")
				if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					return err
				}
				logger.Info("Server stopped")
				return nil
			}, func(err error) {
				logger.Info("Server stopping")
				shutdownCtx, cancelFn := context.WithTimeout(context.Background(), opts.GracePeriod)
				defer cancelFn()
				if err := metricsServer.Shutdown(shutdownCtx); err != nil {
					logger.WithError(err).Error("Server shutdown failed")
					return
				}
				logger.Info("Server stopped gracefully")
			})
		}

		{
			// Set up signal receiver.
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

			g.Add(func() error {
				sig := <-stop
				logrus.WithField("signal", sig.String()).Info("Received signal")
				return nil
			}, func(err error) {
				close(stop)
			})
		}

		// -----------------------------
		// RUN EVERYTHING
		// -----------------------------
		if err := g.Run(); err != nil {
			logrus.WithError(err).Error("Exited unexpectedly")
		}

		logrus.Info("Process exited gracefully")
	},
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
	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		os.Exit(1)
	}
}
