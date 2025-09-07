package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	MySQLDSN    string
	WorkerCount int
}

func Load() *Config {
	port := os.Getenv("PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	db := os.Getenv("MYSQL_DB")
	dsn := user + ":" + pass + "@tcp(" + host + ":3306)/" + db + "?parseTime=true"

	workers, _ := strconv.Atoi(os.Getenv("WORKER_COUNT"))

	return &Config{
		Port:        port,
		MySQLDSN:    dsn,
		WorkerCount: workers,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
