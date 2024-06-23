package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func GetStr(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", errors.New(fmt.Sprintf("nil env key: %s", key))
	}

	return val, nil
}

func GetInt(key string) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return 0, errors.New(fmt.Sprintf("nil env key: %s", key))
	}

	intVal, err := strconv.ParseInt(val, 10, 64)

	return int(intVal), err
}
