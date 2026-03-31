package openshift

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"

	"github.com/openshift/library-go/pkg/crypto"
)

func getRestConfig() (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{}).ClientConfig()
}

type authHandler struct {
	downstream http.Handler
	clientCA   dynamiccertificates.CAContentProvider
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Both authentication and authorization require a client certificate
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		logrus.Info("Client certificate required but not provided")
		http.Error(w, "client certificate required", http.StatusUnauthorized)
		return
	}

	// metricsAllowedClientCommonName is the Common Name (CN) of the client certificate
	// that is authorized to access the metrics endpoint. This corresponds to the
	// well-known Prometheus service account in OpenShift monitoring.
	// See: https://github.com/openshift/enhancements/blob/master/CONVENTIONS.md#metrics
	metricsAllowedClientCommonName := "system:serviceaccount:openshift-monitoring:prometheus-k8s"

	commonName := r.TLS.PeerCertificates[0].Subject.CommonName
	if commonName != metricsAllowedClientCommonName {
		logrus.WithField("commonName", commonName).Info("Access denied")
		http.Error(w, fmt.Sprintf("unauthorized common name: %s", commonName), http.StatusForbidden)
		return
	}

	logrus.WithField("commonName", commonName).Debug("Access granted")
	a.downstream.ServeHTTP(w, r)
}

// MetricsMTLSOptions returns tls.Config that is used to start an HTTP server hosting the Prometheus endpoint
// The http handler does authorization before passing the request to downstream.
// Ref. https://github.com/rhobs/observability-operator/blob/cf377a3068413ade9b7de5452afe18a809b614e8/pkg/operator/operator.go#L146
func MetricsMTLSOptions(ctx context.Context, cert, key string, downstream http.Handler) (*tls.Config, http.Handler, error) {
	restConfig, err := getRestConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get rest config: %w", err)
	}

	var (
		clientCAController    *dynamiccertificates.ConfigMapCAController
		servingCertController *dynamiccertificates.DynamicServingCertificateController
	)

	if cert == "" || key == "" {
		return nil, nil, fmt.Errorf("cert and key are required")
	}

	// DynamicCertKeyPairContent automatically reloads the certificate and key from disk.
	certKeyProvider, err := dynamiccertificates.NewDynamicServingContentFromFiles("serving-cert", cert, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cert-key provider: %w", err)
	}
	if err := certKeyProvider.RunOnce(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize cert/key content: %w", err)
	}

	go certKeyProvider.Run(ctx, 1)

	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize kubernetes client: %w", err)
	}

	clientCAController, err = dynamiccertificates.NewDynamicCAFromConfigMapController(
		"client-ca",
		metav1.NamespaceSystem,
		"extension-apiserver-authentication",
		"client-ca-file",
		kubeClient,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client CA controller: %w", err)
	}

	if err := clientCAController.RunOnce(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to load CA bundle: %w", err)
	}

	go clientCAController.Run(ctx, 1)

	// Only log the events emitted by the certificate controller for now
	// because the controller generates invalid events rejected by the
	// Kubernetes API when used with DynamicServingContentFromFiles.
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(func(format string, args ...interface{}) {
		logrus.Info(fmt.Sprintf(format, args...))
	})

	servingCertController = dynamiccertificates.NewDynamicServingCertificateController(
		&tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
		clientCAController,
		certKeyProvider,
		nil,
		record.NewEventRecorderAdapter(
			eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "observability-operator"}),
		),
	)
	if err := servingCertController.RunOnce(); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize serving certificate controller: %w", err)
	}

	clientCAController.AddListener(servingCertController)
	certKeyProvider.AddListener(servingCertController)

	go servingCertController.Run(1, ctx.Done())

	return crypto.SecureTLSConfig(&tls.Config{
		GetConfigForClient: func(clientHello *tls.ClientHelloInfo) (*tls.Config, error) {
			config, err := servingCertController.GetConfigForClient(clientHello)
			if err != nil {
				return nil, err
			}
			if config == nil {
				// To ensure we rather safely fail connections when the desired config is nil. Safety over availability.
				err := fmt.Errorf("serving certificate controller returned nil TLS configuration")
				return nil, err
			}
			return config, nil
		},
	}), &authHandler{downstream: downstream, clientCA: clientCAController}, nil
}

// MetricsTLSOptions returns tls.Config that is used to start an HTTP server hosting the Prometheus endpoint
// Ref. https://github.com/rhobs/observability-operator/blob/cf377a3068413ade9b7de5452afe18a809b614e8/pkg/operator/operator.go#L146
func MetricsTLSOptions(ctx context.Context, cert, key string) (*tls.Config, error) {

	if cert == "" || key == "" {
		return nil, fmt.Errorf("cert and key are required")
	}

	// DynamicCertKeyPairContent automatically reloads the certificate and key from disk.
	certKeyProvider, err := dynamiccertificates.NewDynamicServingContentFromFiles("serving-cert", cert, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert-key provider: %w", err)
	}
	if err := certKeyProvider.RunOnce(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize cert/key content: %w", err)
	}

	go certKeyProvider.Run(ctx, 1)

	return crypto.SecureTLSConfig(&tls.Config{
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			certPEM, keyPEM := certKeyProvider.CurrentCertKeyContent()

			if certPEM == nil || keyPEM == nil {
				return nil, fmt.Errorf("certificate not ready")
			}

			cert, err := tls.X509KeyPair(certPEM, keyPEM)
			if err != nil {
				return nil, fmt.Errorf("failed to parse key pair: %w", err)
			}

			return &cert, nil
		},
	}), nil
}
