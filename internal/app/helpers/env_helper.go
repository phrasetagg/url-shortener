package helpers

import (
	"os"
)

func GetEnv(envName string, defaultValue string) string {
	env := os.Getenv(envName)

	if env != "" {
		return env
	}

	return defaultValue
}
