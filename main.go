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

	resources, err := aws.Scan(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to scan resources: %v", err)
	}

	for _, resource := range resources {
		println(resource.Region, resource.Service, resource.TypeName, resource.Id)
	}
}
