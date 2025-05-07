package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          int
	Database_Name string
	MongoDBURI    string
	GinMode       string
	JwtSecret     string
}

var Env Config

func LoadConfig() {
	godotenv.Load()

	// Load environment variables
	Env = Config{
		Port:          getEnvAsInt("PORT", 8080),
		Database_Name: getEnv("DATABASE_NAME", ""),
		MongoDBURI:    getEnv("MONGODB_URI", ""),
		GinMode:       getEnv("GIN_MODE", "release"),
		JwtSecret:     getEnv("JWT_SECRET_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if defaultValue == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, strconv.Itoa(defaultValue))
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
