package config

import "os"

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Port      string
	StaticDir string
}

type DBConfig struct {
	Driver string
	DSN    string
}

func Load() *Config {
	staticDir := getEnv("STATIC_DIR", "web/dist")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		staticDir = ""
	}
	return &Config{
		Server: ServerConfig{
			Port:      getEnv("SERVER_PORT", "8080"),
			StaticDir: staticDir,
		},
		DB: DBConfig{
			Driver: getEnv("DB_DRIVER", "sqlite"),
			DSN:    getEnv("DB_DSN", "humq.db"),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
