package utils

import (
	"os"
)

func EnvDefault(env, defaultValue string) string {
	if value, ok := os.LookupEnv(env); ok {
		return value
	}
	return defaultValue
}
