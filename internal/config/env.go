package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	if !CheckEnvs() {
		os.Exit(1)
	} else {
		log.Println("All required environment variables are set")
	}

}

func CheckEnvs() bool {
	requiredEnvs := []string{
		"DISCOGO_HTTP_HOST",
		"DISCOGO_HTTP_PORT",
		"DISCOGO_VERSION",
		"DISCOGO_NAME",
		"DISCOGO_VERSION_NAME",

		"REDIS_HOST",
		"REDIS_PORT",
		"REDIS_PASSWORD",
		"REDIS_DB",

		"HEALTH_CHECK_INTERVAL",
		"REPORT_TOLERANCE_COUNT",
	}

	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatal("Required environment variable not set: " + env)
			return false
		}
	}
	return true
}

func GetRedisServerAddr() string {
	if os.Getenv("REDIS_HOST") == "" || os.Getenv("REDIS_PORT") == "" {
		return "127.0.0.1:6379"
	}
	return os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
}

func GetRedisCredentials() (string, string, int) {
	if os.Getenv("REDIS_HOST") == "" || os.Getenv("REDIS_PORT") == "" || os.Getenv("REDIS_PASSWORD") == "" {
		return "127.0.0.1:6379", "default_pass", 0
	}

	return os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_PASSWORD"),
		getEnvAsInt("REDIS_DB", 0)
}

func GetRedisDB() int {
	return getEnvAsInt("REDIS_DB", 0)
}

func GetRedisPassword() string {
	val := os.Getenv("REDIS_PASSWORD")
	if val == "" {
		return "default_password"
	}
	return val
}

func GetDiscoGoHTTPAddr() string {
	if os.Getenv("DISCOGO_HTTP_HOST") == "" || os.Getenv("DISCOGO_HTTP_PORT") == "" {
		return "127.0.0.1:8080"
	}
	return os.Getenv("DISCOGO_HTTP_HOST") + ":" + os.Getenv("DISCOGO_HTTP_PORT")
}

func GetDiscoGoVersion() string {
	val := os.Getenv("DISCOGO_VERSION")
	if val == "" {
		return "1.0.0"
	}
	return val
}

func GetDiscoGoName() string {
	val := os.Getenv("DISCOGO_NAME")
	if val == "" {
		return "discoGo"
	}
	return val
}

func GetDiscoGoVersionName() string {
	val := os.Getenv("DISCOGO_VERSION_NAME")
	if val == "" {
		return "Artemis"
	}
	return val
}
func IsColorEnabled() bool {
	val := os.Getenv("DISCOGO_LOG_COLOR")
	if val == "" {
		return false
	}
	return val == "true" || val == "1"
}

func GetHealthCheckInterval() int {
	return getEnvAsInt("HEALTH_CHECK_INTERVAL", 30)
}

func GetReportToleranceCount() int64 {
	valStr := os.Getenv("REPORT_TOLERANCE_COUNT")
	if valStr == "" {
		return 5 // Default value if not set
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 5
	}
	return val
}

// getEnvAsInt retrieves an environment variable as an int, or returns the default value if not set or invalid.
func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
