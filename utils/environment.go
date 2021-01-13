package utils

import "os"

var (
	keyENV      = "ENV"
	devENV      = "development"
	defaultPort = "3000"
)

// GetENV ...
func GetENV() string {
	env := os.Getenv(keyENV)

	if len(env) == 0 {
		env = devENV
	}

	return env
}

// GetPort ...
func GetPort() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = defaultPort
	}

	return port
}
