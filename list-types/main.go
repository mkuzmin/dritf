package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cfschema "github.com/hashicorp/aws-cloudformation-resource-schema-sdk-go"
	"log"
	"os"
	"time"
)

var ctx = context.Background()

func main() {

	file, err := os.Create("types.txt")
	if err != nil {
		log.Fatal("Cannot create file: ", err)
	}
	defer file.Close()

	cfTypes := getCfTypes("eu-central-1") // TODO iterate over all regions
	for _, cfType := range cfTypes {
		_, err := file.WriteString(cfType + "\n")
		if err != nil {
			log.Fatal("Cannot write to file: ", err)
		}
	}
}

func getCfTypes(region string) []string {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	cfg.Region = region

	var cfTypes []string
	println("\n=== Mutable")
	cfTypes = append(cfTypes, getCfTypesProvisioning(cfg, types.ProvisioningTypeFullyMutable)...)

	println("\n=== Immutable")
	cfTypes = append(cfTypes, getCfTypesProvisioning(cfg, types.ProvisioningTypeImmutable)...)

	println("\n=== Non-provisionable")
	cfTypes = append(cfTypes, getCfTypesProvisioning(cfg, types.ProvisioningTypeNonProvisionable)...)

	return cfTypes
}

func getCfTypesProvisioning(cfg aws.Config, provisioningType types.ProvisioningType) []string {
	cf := cloudformation.NewFromConfig(cfg)
	input := cloudformation.ListTypesInput{
		Type:             types.RegistryTypeResource,
		Visibility:       types.VisibilityPublic,
		DeprecatedStatus: types.DeprecatedStatusLive,
		Filters: &types.TypeFilters{
			Category: types.CategoryAwsTypes,
		},
		ProvisioningType: provisioningType,
	}

	var cfTypes []string
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

			print("\r" + *ts.TypeName + "                                                     ")

			if isListable(cf, *ts.TypeName) {
				cfTypes = append(cfTypes, *ts.TypeName)
			}
		}
	}

	return cfTypes
}

func isListable(cf *cloudformation.Client, cfType string) bool {
	typeInput := cloudformation.DescribeTypeInput{
		Type:     types.RegistryTypeResource,
		TypeName: &cfType,
	}

	time.Sleep(500 * time.Millisecond) // prevent API throttling
	describeType, err := cf.DescribeType(ctx, &typeInput)
	if err != nil {
		log.Fatal("Cannot describe type", err)
	}
	if describeType.Schema == nil {
		log.Fatal("schema is nil")
	}

	var r cfschema.Resource
	err = json.Unmarshal([]byte(*describeType.Schema), &r)
	if err != nil {
		log.Fatal("Error parsing JSON: ", err)
	}

	if _, exists := r.Handlers["list"]; exists {
		return true
	}

	println()
	return false
}
