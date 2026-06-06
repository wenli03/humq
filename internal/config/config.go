package config

import "os"

type Config struct {
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
	Kafka  KafkaConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	DSN string
}

type JWTConfig struct {
	Secret             string
	ExpireHours        int
	RefreshExpireHours int
}

type KafkaConfig struct {
	BootstrapServers string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{Port: getEnv("SERVER_PORT", "8080")},
		DB:     DBConfig{DSN: getEnv("DB_DSN", "host=localhost user=humq password=humq dbname=humq port=5432 sslmode=disable")},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "humq-secret-change-in-production"),
			ExpireHours:        24,
			RefreshExpireHours: 168,
		},
		Kafka: KafkaConfig{BootstrapServers: getEnv("KAFKA_BOOTSTRAP", "localhost:9092")},
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
