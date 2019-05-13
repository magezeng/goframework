package Models

type DBConfig struct {
	Name       string `yaml:"name" json:"name"`
	ConnectStr string `yaml:"connectStr" json:"connectStr"`
}

type HTTPConfig struct {
	MaxIdleConns        int `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxIdleConnsPerHost int `yaml:"maxIdleConnsPerHost" json:"maxIdleConnsPerHost"`
	IdleConnTimeout     int `yaml:"idleConnTimeout" json:"idleConnTimeout"`
}

type Upload struct {
	Path          string `yaml:"path" json:"path"`
	MaxFileSizeMB int    `yaml:"maxFileSizeMB" json:"maxFileSizeMB"`
}

type Config struct {
	DB     []DBConfig `yaml:"db" json:"db"`
	HTTP   HTTPConfig `yaml:"http" json:"http"`
	Upload Upload     `yaml:"upload" json:"upload"`
}
