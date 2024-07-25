package config

import (
	"log"
	"strconv"
)

func toInt64(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

func setDefaults() {
	if Token == "" {
		log.Fatal("Token required")
	}

	if DatabaseURI == "" {
		log.Fatal("DatabaseURI required")
	}

	if OwnerId == 0 {
		OwnerId = 5938660179
	}

}
