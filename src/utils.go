package main

import (
	"fmt"
	"log/slog"
	"os"
)

/*
ReadRequiredEnvVar reads a specified environment variable and returns the value,
or exits with status 1 if the value is unset or empty.
*/
func ReadRequiredEnvVar(name string) string {
	value := os.Getenv(name)
	if value == "" {
		slog.Error(fmt.Sprintf("missing required environment variable %s", name))
		os.Exit(1)
	}
	return value
}

/*
ReadEnvVarWithDefault reads a specified environment variable and returns the value,
or returns a specified default value if the value is unset or empty.
*/
func ReadEnvVarWithDefault(name string, defaultVal string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultVal
	}
	return value
}
