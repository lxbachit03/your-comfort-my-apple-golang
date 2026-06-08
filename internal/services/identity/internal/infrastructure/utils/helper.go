package utils

import (
	"log"
	"os"
)

func GetWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("❌ Unable to get working dir:", err)
	}
	return dir
}
