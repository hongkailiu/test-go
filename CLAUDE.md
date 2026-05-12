# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go implementation of the Cincinnati update graph server for OpenShift. Cincinnati is the update protocol that serves upgrade paths to OpenShift clusters. The server:
- Builds and serves update graphs showing valid upgrade paths between OpenShift versions
- Scrapes container registries for release metadata
- Supports multi-arch releases (amd64, arm64, s390x, ppc64le, multi)
- Provides HTTP endpoints for graph queries and metrics

## Development Commands

### Running Locally
```bash
# Run the server with mock data
MOCK_DIR=mock go run ./cmd/cincinnati --log-level debug --max-concurrency 1 --metrics-tls-disabled

# Run with real data (requires registry access)
go run ./cmd/cincinnati --registry https://quay.io --repo openshift-release-dev/ocp-release
```

### Testing
```bash
# Run unit tests
make unit

# Run all tests (same as unit)
make test

# Run integration tests (requires KUBECONFIG set to a real cluster)
make integration-test

# Run a specific test
go test -v ./pkg/cincinnati -run TestGraphBuilder
```

### Code Quality
```bash
# Run linter (golangci-lint)
make lint

# Format imports
make imports

# Format Go code and tidy dependencies
make generate-go

# Run all verification (tests, imports, formatting, lint, check for uncommitted changes)
make verify

# YAML linting
make yaml-lint

# Full verification (includes YAML)
make verify-all
```

### Building
```bash
# Build container image
podman build -t cincinnati-go -f images/cincinnati/Containerfile .
```

## Architecture

### Core Components

**cmd/cincinnati/main.go**: Application entry point that:
- Initializes HTTP servers (main + metrics)
- Sets up periodic workers using oklog/run groups
- Configures GraphBuilder and Repo components
- Handles graceful shutdown

**pkg/cincinnati/builder.go (GraphBuilder)**: Orchestrates graph construction:
- Caches graph data from local directories (refreshed every 1 minute)
- Builds the full upgrade graph (refreshed every 2 hours)
- Applies graph handlers sequentially: repo scraping → graph data shaping
- Loads graphs from disk on startup, falls back to building from scratch

**pkg/cincinnati/repo.go (Repo)**: Registry scraper:
- Fetches container image tags and metadata from registries (Quay.io)
- Extracts version information from release images
- Converts registry tags to graph nodes and edges
- Supports concurrent scraping with configurable max-concurrency

**pkg/cincinnati/graph.go (Graph)**: Graph data structures and filtering:
- Graph contains Nodes (releases), Edges (unconditional upgrades), and ConditionalEdges (conditional upgrades with risks)
- Supports per-request filtering by channel, architecture, and version
- Handles multi-arch releases and architecture-specific upgrade paths

**pkg/cincinnati/server.go**: HTTP handlers:
- `/api/upgrades_info/graph` and `/api/upgrades_info/v1/graph`: Query upgrade graphs by channel, arch, version
- `/api/upgrades_info/graph-data`: Download graph data as tar.gz
- `/healthz` and `/readyz`: Health checks
- `/version`: Version information

**pkg/cincinnati/data.go (CincinnatiGraphData)**: Loads graph metadata from disk:
- Reads blocked-edges, channels, and removed-edges from `graph-data/` directory
- Shapes the graph by applying channel memberships, blocking rules, and conditional risks

### Data Flow

1. **Initialization**: Load graph from disk file (data/graph.json.gz) if available
2. **CacheGraphData** (every 1 minute): Load channel info, blocked edges, etc. from `graph-data/` directory
3. **CacheGraph** (every 2 hours):
   - Repo scrapes registry tags and converts them to nodes/edges
   - CincinnatiGraphData applies channel memberships and blocking rules
   - Result cached in memory and written to disk
4. **Client request**: GraphBuilder filters cached graph by channel/arch/version and returns matching subgraph

### Directory Structure

- `data/`: Version metadata organized by major.minor version (e.g., `data/4.15/4.15.0-x86_64.yaml`)
- `mock/`: Mock data for local development without registry access
- `pkg/cincinnati/`: Core graph building and serving logic
- `pkg/openshift/`: OpenShift-specific utilities (e.g., metrics mTLS configuration)
- `pkg/util/`: Utility functions for gzip, tar operations
- `cmd/cincinnati/`: Main application entry point

## Testing with Metrics Endpoint

The metrics endpoint supports both HTTP and HTTPS (with mTLS). From README.md:

### mTLS (requires client certificates)
```bash
# Test from within OpenShift cluster with proper client certs
oc --kubeconfig /tmp/ota-stage.c -n openshift-monitoring exec -c prometheus pod/prometheus-k8s-0 -- \
  curl -si --cacert /etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt \
  --cert /etc/prometheus/secrets/metrics-client-certs/tls.crt \
  --key /etc/prometheus/secrets/metrics-client-certs/tls.key \
  https://cincinnati.cincinnati-go.svc.cluster.local:9090/metrics
```

Note: ServiceMonitor in User Workload Monitoring cannot use `caFile: /etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt` or `scrapeClass: tls-client-certificate-auth` - these only work in Cluster Monitoring. For mTLS to work, generate cert/key files for a ServiceAccount signed by Kubernetes.

### HTTP (for local development)
```bash
# Run with metrics TLS disabled
go run ./cmd/cincinnati --metrics-tls-disabled
curl http://localhost:9090/metrics
```

## Important Implementation Details

- **Architecture handling**: Tags use suffixes like `-x86_64`, `-aarch64`, `-s390x`, `-ppc64le`, `-multi`. The graph.go `getArch()` function defaults to amd64 if no suffix is found.
- **Multi-arch support**: Multi-arch releases (GA since 4.13, tech preview in 4.11-4.12) have special handling. Cincinnati only shows multi-arch upgrade paths from 4.12+.
- **Semantic versioning**: Uses `github.com/blang/semver/v4` for version parsing. Note: Use `semver.Parse()` for consistency (there's a TODO about this in the code).
- **Caching**: Uses `github.com/patrickmn/go-cache` for in-memory caching with expiration
- **Graceful shutdown**: Uses `github.com/oklog/run` for coordinated actor shutdown
- **Logging**: Uses logrus with structured logging; set log level with `--log-level` flag

## Common Gotchas

- Integration tests (`TestIntegration*`) require `TEST_INTEGRATION=1` environment variable and a valid KUBECONFIG
- The `--mock-dir` flag automatically sets `--graph-data-dir` to `$MOCK_DIR/graph-data`
- The linter runs with `--new-from-rev 6d2c3ff3a6ca6ede20fee99f9db36cd3992349e5` to only check new changes
- When working with graph filtering, remember that ConditionalEdges are not arch-specific and persist across architectures
