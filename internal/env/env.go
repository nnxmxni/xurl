package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func GetInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		valAsInt, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}

		return valAsInt
	}

	return fallback
}
