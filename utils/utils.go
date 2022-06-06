package utils

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/google/uuid"
)

const MaxInt = int(^uint(0) >> 1)

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

func SetPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		return "8080"
	}

	return port
}

func CalculateDistanceKM(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}
	fmt.Println(dist)

	return dist
}
