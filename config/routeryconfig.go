package config

type RouteryConfig struct {
	Logging struct {
		File    bool
		Path    string
		Verbose bool
	}

	Frontend []FrontendConfig
	Docker   []DockerConfig
	Auth []AuthConfig
}

type DockerConfig struct {
	IP   string
	Port int
	SSL  bool
	CA   string `json:"ca,omitempty"`
	Cert string `json:"cert,omitempty"`
	Key  string `json:"key,omitempty"`
}

type FrontendConfig struct {
	Hostname string
	Port     int
	SSL      bool
	Cert     string `json:"cert,omitempty"`
	Key      string `json:"key,omitempty"`
	CA       string `json:"ca,omitempty"`
}

type AuthConfig struct {
	Type string
	Hostname string
	Port int
	Arguments string  `json:"arguments,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
