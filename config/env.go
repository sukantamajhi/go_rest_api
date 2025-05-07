package config

import (
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
		MongoDBURI:    getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		GinMode:       getEnv("GIN_MODE", "release"),
		JwtSecret:     getEnv("JWT_SECRET_KEY", "secret"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
