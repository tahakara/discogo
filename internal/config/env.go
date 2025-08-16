package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
func GetMemcachedServerAddr() (string, error) {
	host := os.Getenv("MEMCACHED_HOST")
	port := os.Getenv("MEMCACHED_PORT")
	if host == "" || port == "" {
		return "", errors.New("MEMCACHED_HOST or MEMCACHED_PORT environment variable not set")
	}
	return host + ":" + port, nil
}

func GetDiscoGoHTTPAddr() (string, error) {
	host := os.Getenv("DISCOGO_HTTP_HOST")
	port := os.Getenv("DISCOGO_HTTP_PORT")
	if host == "" || port == "" {
		return "", errors.New("DISCOGO_HTTP_HOST or DISCOGO_HTTP_PORT environment variable not set")
	}
	return host + ":" + port, nil
}

func GetDiscoGoVersion() (string, error) {
	version := os.Getenv("DISCOGO_VERSION")
	if version == "" {
		return "", errors.New("DISCOGO_VERSION environment variable not set")
	}
	return version, nil
}

func GetDiscoGoName() (string, error) {
	name := os.Getenv("DISCOGO_NAME")
	if name == "" {
		return "", errors.New("DISCOGO_NAME environment variable not set")
	}
	return name, nil
}

func GetDiscoGoVersionName() (string, error) {
	versionName := os.Getenv("DISCOGO_VERSION_NAME")
	if versionName == "" {
		return "", errors.New("DISCOGO_VERSION_NAME environment variable not set")
	}
	return versionName, nil
}

func IsColorEnabled() bool {
	val := os.Getenv("DISCOGO_LOG_COLOR")
	if val == "" {
		return false
	}
	return val == "true"
}
