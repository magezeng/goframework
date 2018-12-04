package do

// Configuration 系统配置的结构，一个文件包含所有
type Configuration struct {
	Db struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Database string `json:"database"`
		Password string `json:"password"`
		Port     int    `json:"port"`
	}
	Logger struct {
		FileExt    string `json:"fileExt"`
		FileName   string `json:"fileName"`
		SavePath   string `json:"savePath"`
		TimeFormat string `json:"timeFormat"`
	}
}
