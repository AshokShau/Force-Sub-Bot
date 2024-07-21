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

	if LoggerId == 0 {
		log.Fatal("LoggerId required")
	}

	if OwnerId == 0 {
		OwnerId = 5938660179
	}

	if DbName == "" {
		DbName = "ForceSub"
	}

}

// FindInInt64Slice Find takes a slice and looks for an element in it. If found it will
// return true, otherwise it will return a bool of false.
func FindInInt64Slice(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// RemoveFromInt64Slice Find takes a slice and looks for an element in it. If found it will
func RemoveFromInt64Slice(s []int64, r int64) []int64 {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
