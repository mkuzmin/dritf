package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/aws/aws-sdk-go-v2/service/account/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"log"
	"slices"
)

func Scan(ctx context.Context, cfg *Config) chan Result {
	resultChan := make(chan Result)
	go listResources(ctx, cfg, resultChan)
	return resultChan
}

func listResources(ctx context.Context, cfg *Config, resultChan chan Result) {
	awsConfig, err := config.LoadDefaultConfig(
		ctx,
		config.WithAppID("github.com/mkuzmin/dritf"),
	)
	if err != nil {
		log.Fatalf("failed to create AWS config: %v", err)
	}

	enabledRegions := getEnabledRegions(ctx, awsConfig)

	for _, region := range cfg.Regions {
		if !slices.Contains(enabledRegions, region) {
			continue
		}
		log.Println(region)

		client := cloudcontrol.NewFromConfig(
			awsConfig,
			func(o *cloudcontrol.Options) { o.Region = region },
		)

		for _, service := range cfg.Services {
			for _, resourceType := range service.ResourceTypes {
				name := fmt.Sprintf("AWS::%s::%s", service.Name, resourceType.Name)
				input := cloudcontrol.ListResourcesInput{TypeName: &name}

				paginator := cloudcontrol.NewListResourcesPaginator(client, &input)
				for paginator.HasMorePages() {
					output, err := paginator.NextPage(ctx)
					if err != nil {
						resultChan <- Result{
							Error: fmt.Errorf("failed to list resources for '%s' (%s): %v", name, region, err),
						}
						break
					}

					for _, res := range output.ResourceDescriptions {
						id := *res.Identifier
						resultChan <- Result{
							Resource: Resource{
								Region:   region,
								Service:  service.Name,
								TypeName: resourceType.Name,
								Id:       id,
							},
						}
					}
				}
			}
		}
	}

	close(resultChan)
}

func getEnabledRegions(ctx context.Context, awsConfig aws.Config) []string {
	accountClient := account.NewFromConfig(awsConfig)

	paginator := account.NewListRegionsPaginator(accountClient, &account.ListRegionsInput{})

	var enabledRegions []string
	for paginator.HasMorePages() {
		listRegionsOutput, err := paginator.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to list regions: %v", err)
		}

		for _, regionDetail := range listRegionsOutput.Regions {
			if regionDetail.RegionOptStatus != types.RegionOptStatusEnabledByDefault &&
				regionDetail.RegionOptStatus != types.RegionOptStatusEnabled {
				continue
			}
			if regionDetail.RegionName == nil {
				log.Fatalf("region has no name: %v", err)
			}
			enabledRegions = append(enabledRegions, *regionDetail.RegionName)
		}
	}
	return enabledRegions
}
