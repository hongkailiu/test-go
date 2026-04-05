# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go implementation of the Cincinnati update graph server for OpenShift. The server provides upgrade path information to OpenShift clusters by:
- Fetching release metadata from container registries (Quay.io)
- Loading graph data (channels, blocked edges, raw metadata) from disk
- Building and caching upgrade graphs
- Serving filtered graphs based on channel, architecture, and current version

## Development Commands

### Running Locally
```bash
# Run with mock data (fastest for development)
MOCK_DIR=mock go run ./cmd/cincinnati --log-level debug --max-concurrency 1 --metrics-tls-disabled

# Run against real registry (slower, requires network access)
go run ./cmd/cincinnati --log-level debug --registry https://quay.io --repo openshift-release-dev/ocp-release
```

### Building
```bash
# Build binary to _out/cincinnati
hack/build.sh

# Build container image
podman build -t test001 -f images/cincinnati/Containerfile .
```

### Testing
```bash
# Run unit tests (uses gotestsum)
make test

# Run specific test
go test -v ./pkg/cincinnati -run TestBuildGraph

# Run integration tests (requires OpenShift cluster via KUBECONFIG)
make integration-test
```

### Code Quality
```bash
# Full verification (runs tests, linting, formatting, and checks for uncommitted changes)
make verify

# Individual steps:
make lint           # golangci-lint (configured in .golangci.yml)
make imports        # Organize imports with gci
make generate-go    # Format and tidy go.mod
make yaml-lint      # Lint YAML files

# Verify everything including YAML
make verify-all
```

## Architecture

### Core Components

**GraphBuilder** (`pkg/cincinnati/builder.go`):
- Orchestrates graph construction through a pipeline of handlers
- Loads graph data from disk and caches it
- Periodically rebuilds the graph (every 2 hours for full graph, every 1 minute for graph data)
- Two key cache entries: `openshift-upgrade-graph` and `cincinnati-graph-data`
- Server is "ready" only when both cache entries are populated

**Repo** (`pkg/cincinnati/repo.go`):
- Fetches release metadata from container registries
- Converts registry tags into graph nodes and edges
- Handles concurrent registry scraping (controlled by `--max-concurrency`)
- Can use mock data when `MOCK_DIR` is set

**Server** (`pkg/cincinnati/server.go`):
- Exposes HTTP API endpoints via gin framework
- Main server (default :8080) serves the upgrade graph API
- Metrics server (default :9090) serves Prometheus metrics with optional mTLS

**Graph** (`pkg/cincinnati/graph.go`):
- Core data structure representing upgrade paths
- Nodes: individual releases with metadata
- Edges: allowed upgrade paths
- ConditionalEdges: upgrade paths with conditions
- Supports filtering by channel, architecture, version, and ID

### Startup Flow

1. Initial work (blocking): `CacheGraphData()` → `CacheGraph()`
2. Background workers start:
   - Graph data refresh: every 1 minute
   - Full graph rebuild: every 2 hours
3. HTTP servers start (main + metrics)
4. Server becomes "ready" once initial caching completes

### Data Directory Structure

The `data/` directory contains versioned subdirectories (4.0, 4.1, 4.10, etc.) with release information files. When using mock mode (`MOCK_DIR=mock`), the server loads data from the mock directory instead.

Graph data directory structure:
- `blocked-edges/`: Edge blocks (prevent specific upgrades)
- `channels/`: Channel definitions
- `raw/`: Raw metadata files
- `LICENSE`, `version`: Metadata files

## API Endpoints

All endpoints use the `/api/upgrades_info` prefix:

- `GET /api/upgrades_info/graph?channel=<channel>&arch=<arch>&version=<version>&id=<id>`
  - `channel`: required (e.g., "stable-4.16")
  - `arch`: optional, defaults to "amd64"
  - `version`: optional semver filter
  - `id`: optional release ID filter

- `GET /api/upgrades_info/graph-data`: Downloads graph data as tar.gz

Health checks:
- `GET /healthz`: Always returns 200 if server is running
- `GET /readyz`: Returns 200 only when cache is populated (server is ready)
- `GET /version`: Server version information

Metrics:
- `GET /metrics` (on metrics port): Prometheus metrics, supports mTLS

## Configuration

Environment variables:
- `CINCINNATI_REGISTRY`: Override registry URL (default: https://quay.io)
- `CINCINNATI_REPO`: Override repo path (default: openshift-release-dev/ocp-release)
- `MOCK_DIR`: Use mock data instead of real registry

Command-line flags: see `cmd/cincinnati/main.go` init() or run with `--help`

## Testing Notes

- Unit tests use golden files (via cupaloy) in `pkg/cincinnati/testdata/`
- Integration tests require a real OpenShift cluster (set `TEST_INTEGRATION=1`)
- Mock mode is recommended for fast iteration during development
- The `--max-concurrency 1` flag is useful for debugging to avoid parallel goroutines

## Import Organization

This project uses gci to organize imports in the following order:
1. Standard library
2. Default (third-party)
3. k8s.io packages
4. github.com/openshift packages
5. Local module imports

Run `make imports` to auto-format.
