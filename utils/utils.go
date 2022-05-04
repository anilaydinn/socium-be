package utils

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID(length int) string {
	uuid := uuid.New().String()

	uuid = strings.ReplaceAll(uuid, "-", "")

	if length < 1 {
		return uuid
	}
	if length > len(uuid) {
		length = len(uuid)
	}

	return uuid[0:length]
}

func GetDBUrl() string {
	databaseUrl := os.Getenv("DATABASE_URL")

	if databaseUrl == "" {
		return "mongodb://localhost:27017"
	}

	return databaseUrl
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
func RemoveElement(s []string, id string) []string {
	index := linearSearch(s, id)

	if index != -1 {
		return append(s[:index], s[index+1:]...)
	} else {
		return s
	}
}

func linearSearch(s []string, id string) int {
	for i, n := range s {
		if n == id {
			return i
		}
	}
	return -1
}
