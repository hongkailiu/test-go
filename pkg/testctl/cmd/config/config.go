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
