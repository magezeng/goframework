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

type Config struct {
	DB     []DBConfig   `yaml:"db" json:"db"`
	HTTP   HTTPConfig   `yaml:"http" json:"http"`
}
