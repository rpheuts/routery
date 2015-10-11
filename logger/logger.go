package logger

import (
	"log"
	"os"
)

func SetLogging(file bool, path string) {
	if !file {
		return
	}

	f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
}
