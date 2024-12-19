package main

import (
	"log"
	"os"
)

type Config struct {
	ClientID     string
	ClientSecret string
	BucketName   string
}

// LoadConfig loads configuration values from environment variables.
func LoadConfig() Config {
	return Config{
		ClientID:     mustGetEnv("SFMC_CLIENT_ID"),
		ClientSecret: mustGetEnv("SFMC_CLIENT_SECRET"),
		BucketName:   mustGetEnv("AWS_S3_BUCKET_NAME"),
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
