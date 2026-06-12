package utils

import (
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
