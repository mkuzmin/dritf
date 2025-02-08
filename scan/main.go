package main

import (
	"context"
	"github.com/mkuzmin/dritf/common"
	"log"
)

var ctx = context.Background()

func main() {
	cfg, err := common.LoadConfig("mapping.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	common.ScanTypes(cfg, ctx)
}
