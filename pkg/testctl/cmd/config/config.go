package config

// Config ...
type Config struct {
	Verbose bool
}

type HttpConfig struct {
	Config
	PProf   bool
	Version string
}
