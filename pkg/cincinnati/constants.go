package cincinnati

const (
	// DefaultMaxConcurrency is the value for max concurrency to fetch images from the repository.
	DefaultMaxConcurrency = 3

	DefaultPort        = ":8080"
	DefaultMetricsPort = ":9090"

	MetadataKeyChannels     = "io.openshift.upgrades.graph.release.channels"
	MetadataKeyArchitecture = "release.openshift.io/architecture"
)
