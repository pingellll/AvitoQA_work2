package utils

import (
	"os"
	"strings"
)

func GetBaseURL() string {
	base := strings.TrimSpace(os.Getenv("API_URL"))
	return strings.TrimRight(base, "/")
}
