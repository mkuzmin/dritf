package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"log"
	"os"
)

var ctx = context.Background()

func main() {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	cfg.Region = "eu-central-1" // TODO iterate over all regions

	cf := cloudformation.NewFromConfig(cfg)
	input := cloudformation.ListTypesInput{
		Type:             types.RegistryTypeResource,
		Visibility:       types.VisibilityPublic,
		DeprecatedStatus: types.DeprecatedStatusLive,
		Filters: &types.TypeFilters{
			Category: types.CategoryAwsTypes,
		},
		ProvisioningType: types.ProvisioningTypeFullyMutable, // TODO read other types
	}

	file, err := os.Create("types.txt")
	if err != nil {
		log.Fatal("Cannot create file: ", err)
	}
	defer file.Close()

	paginator := cloudformation.NewListTypesPaginator(cf, &input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			log.Fatal("Cannot list types: ", err)
		}

		for _, ts := range output.TypeSummaries {
			if ts.TypeName == nil {
				log.Fatal("List type is nil")
			}
			_, err := file.WriteString(*ts.TypeName + "\n")
			if err != nil {
				log.Fatal("Cannot write to file: ", err)
			}
		}
	}
}
