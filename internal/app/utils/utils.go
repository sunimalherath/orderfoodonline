// Package utils: contains utilitarian functions.
package utils

import "os"

func GetEnvVar(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}
