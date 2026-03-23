package cincinnati

const (
	// DefaultMaxConcurrency is the value for max concurrency to fetch images from the repository.
	DefaultMaxConcurrency = 3

	// MetadataKeyManifestRef records the manifest reference of the image
	// Used only internally by Cincinnati Rust implementation while handling graph-data's blocked edges.
	// It is redundant information as it is part of node's version.
	// https://redhat-internal.slack.com/archives/CJ1J9C3V4/p1773883232868309?thread_ts=1773881763.378719&cid=CJ1J9C3V4
	MetadataKeyManifestRef  = "io.openshift.upgrades.graph.release.manifestref"
	MetadataKeyChannels     = "io.openshift.upgrades.graph.release.channels"
	MetadataKeyArchitecture = "release.openshift.io/architecture"
)
