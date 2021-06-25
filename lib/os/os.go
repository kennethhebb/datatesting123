package os

import (
	"os"
	"strconv"
)

func GetStringEnv(key string) string {
	return os.Getenv(key)
}

func GetIntEnv(key string) int {
	val, set := os.LookupEnv(key)
	if !set {
		return 0
	}
	ival, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return ival
}

func Exit(code int) {
	os.Exit(code)
}
