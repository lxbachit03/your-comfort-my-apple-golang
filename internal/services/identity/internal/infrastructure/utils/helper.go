package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
)

func GetWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("❌ Unable to get working dir")
	}
	return dir
}

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
