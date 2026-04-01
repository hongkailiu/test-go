# Build stage
FROM registry.access.redhat.com/ubi9/go-toolset:9.7 AS builder

# Copy source code
COPY --chown=1001:0 . .

# Build the binary
RUN OUT_DIR=/tmp hack/build.sh

# Runtime stage
FROM registry.access.redhat.com/ubi9-minimal:9.7

# Copy the binary from builder stage
COPY --from=builder /tmp/validate-release /usr/local/bin/

# Run the binary
ENTRYPOINT ["/usr/local/bin/validate-release"]