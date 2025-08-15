package config

// AppConfig 應用配置結構
type AppConfig struct {
	ServerPort string
	Database   *DatabaseConfig
}

// NewAppConfig 創建新的應用配置
func NewAppConfig() *AppConfig {
	return &AppConfig{
		ServerPort: "8080",
		Database:   NewDatabaseConfig(),
	}
}
