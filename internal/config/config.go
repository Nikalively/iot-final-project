package config

import "os"

type Config struct {
	RedisAddr string
	Port      string
}

func LoadConfig() *Config {
	return &Config{
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		Port:      getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
