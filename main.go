package main

import (
	"duplicates-finder/gcp"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	// Set up a zap logger
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "message",
	}

	// Set up a zap logger
	config := zap.Config{
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"./duplicates.log"},
		ErrorOutputPaths: []string{"./duplicates_error.log"},
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	// Replace Global logger
	zap.ReplaceGlobals(logger)
}

func main() {
	log.Println("Creating GCP client...")
	client, err := gcp.NewClient()
	if err != nil || client == nil {
		log.Println("Error creating GCP client", err)
	}

	log.Println("Starting to find duplicates...")
	client.Start()
}
