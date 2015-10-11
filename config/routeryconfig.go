package config

type RouteryConfig struct {
	Logging struct {
		File    bool
		Path    string
		Verbose bool
	}

	Frontend []FrontendConfig
	Docker   []DockerConfig
}

type DockerConfig struct {
	IP   string
	Port int
	SSL  bool
	CA   string
	Cert string
	Key  string
}

type FrontendConfig struct {
	Hostname string
	Port     int
	SSL      bool
	Cert     string
	Key      string
	CA       string
}
