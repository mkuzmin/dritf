package main

import (
	"context"
	"github.com/mkuzmin/dritf/aws"
	"log"
)

func main() {
	ctx := context.Background()

	cfg, err := aws.LoadConfig("dritf.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	results := aws.Scan(ctx, cfg)

	for _, result := range results {
		if result.Error != nil {
			log.Printf("failed to scan resources: %v", result.Error)
			continue
		}
		res := result.Resource
		println(res.Region, res.Service, res.TypeName, res.Id)
	}
}
