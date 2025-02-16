package main

import (
	"context"
	"fmt"
	"github.com/mkuzmin/dritf/aws"
	"log"
)

func main() {
	ctx := context.Background()

	cfg, err := aws.LoadConfig("dritf.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	resultChan := aws.Scan(ctx, cfg)

	for result := range resultChan {
		if result.Error != nil {
			log.Printf("failed to scan resources: %v", result.Error)
			continue
		}
		res := result.Resource
		fmt.Println(res.Region, res.Service, res.TypeName, res.Id)
	}
}
