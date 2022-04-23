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
