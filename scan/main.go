package main

import (
	"context"
	"log"
)

var ctx = context.Background()

func main() {
	cfg, err := loadConfig("mapping.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	scanTypes(cfg)
}
