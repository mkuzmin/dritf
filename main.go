package main

import (
	"fmt"
	"github.com/mkuzmin/dritf/aws"
	"log"
)

func main() {
	cfg, err := aws.LoadConfig("dritf.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	fmt.Printf("Regions: %v\n", cfg.Regions)
}
