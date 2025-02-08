package main

import (
	"context"
	"fmt"
	"github.com/mkuzmin/dritf/common"
	"log"
)

var ctx = context.Background()

func main() {
	cfg, err := common.LoadConfig("mapping.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	resourceChan := common.ScanTypes(ctx, cfg)
	for res := range resourceChan {
		fmt.Printf("%s,%s,%s,%s,%s\n", res.Account, res.Region, res.Service, res.TypeName, res.Id)
	}
}
