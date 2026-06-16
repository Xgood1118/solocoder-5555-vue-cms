package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      int
	JWTSecret string
	DataDir   string
	UploadDir string
}

var AppConfig *Config

func LoadConfig() {
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-jwt-secret-key-change-in-production"
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}

	AppConfig = &Config{
		Port:      port,
		JWTSecret: jwtSecret,
		DataDir:   dataDir,
		UploadDir: uploadDir,
	}
}
