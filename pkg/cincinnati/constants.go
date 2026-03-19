package cincinnati

const (
	// DefaultMaxConcurrency is the value for max concurrency to fetch images from the repository.
	DefaultMaxConcurrency = 3

	MetadataKeyManifestRef  = "io.openshift.upgrades.graph.release.manifestref"
	MetadataKeyChannels     = "io.openshift.upgrades.graph.release.channels"
	MetadataKeyArchitecture = "release.openshift.io/architecture"
)
