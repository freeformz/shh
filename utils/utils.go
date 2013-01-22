package utils

import (
	"os"
	"strconv"
)

func Ui64toa(val uint64) string {
	return strconv.FormatUint(val, 10)
}

// Checks to see if a path exists or not
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Returns the value of $env from the OS and if it's empty, returns def
func GetEnvWithDefault(env string, def string) string {
  tmp := os.Getenv(env)

  if tmp == "" {
    return def
	}

  return tmp
}

// Returns the value of $env from the OS and if it's empty, returns def
func GetEnvWithDefaultInt(env string, def int) int {
  tmp := os.Getenv(env)

  if tmp == "" {
    return def
	}

  return tmp
}
