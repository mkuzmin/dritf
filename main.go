package main

import (
	"context"
	"fmt"
	"github.com/mkuzmin/dritf/aws"
	"github.com/mkuzmin/dritf/terraform"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	tfDir := "."
	if len(os.Args) > 1 {
		tfDir = os.Args[1]
	}
	tfResources, err := terraform.GetResources(ctx, tfDir)
	if err != nil {
		log.Fatalf("failed to read Terraform state: %v", err)
	}

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

		tfResource := tfResources.FindResource(res.TypeConfig.TfConfig.Name, res.Id)
		var checkbox string
		var tfAddress string
		if tfResource != nil {
			checkbox = "[V]"
			tfAddress = fmt.Sprint("-> ", tfResource.Address)
		} else {
			checkbox = "[ ]"
		}

		fmt.Println(
			checkbox,
			res.Region,
			res.Service,
			res.TypeName,
			res.Id,
			tfAddress,
		)
	}
}
