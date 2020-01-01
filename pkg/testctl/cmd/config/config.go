package config

// Config ...
type Config struct {
	Verbose bool
}

// HttpConfig ...
type HttpConfig struct {
	Config
	PProf   bool
	Version string
}

type ProwConfig struct {
	Config
	KubeConfigPath string
}