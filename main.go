package main

import (
	"context"
	"duplicates-finder/gcp"
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var wg sync.WaitGroup

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
	user_id := "2LGMW74VzpaUJcW0wnh2ZQUwaqn2"

	log.Println("Creating GCP client...")
	client, err := gcp.NewClient()
	if err != nil || client == nil {
		log.Println("Error creating GCP client", err)
	}

	log.Println("Getting user profiles...")
	paths, err := client.GetUserProfiles(user_id)
	if err != nil {
		log.Println("Error getting user profiles", err)
	}

	log.Println("Listing zip objects...", len(paths))
	for _, path := range paths {
		wg.Add(1)
		go client.ListZipObjects(context.Background(), path, &wg)
	}

	wg.Wait()
}
