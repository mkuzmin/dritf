package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/aws/aws-sdk-go-v2/service/account/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"log"
	"slices"
)

func scanTypes(cfg *Config) {
	awsConfig, err := config.LoadDefaultConfig(
		ctx,
		config.WithAppID("github.com/mkuzmin/dritf"),
	)
	if err != nil {
		log.Fatal("failed to create AWS client: ", err)
	}

	enabledRegions := getEnabledRegions(awsConfig)

	for _, region := range cfg.Regions {
		if !slices.Contains(enabledRegions, region) {
			continue
		}

		ccClient := cloudcontrol.NewFromConfig(
			awsConfig,
			func(o *cloudcontrol.Options) { o.Region = region },
		)

		for _, svc := range cfg.Services {
			for _, rt := range svc.ResourceTypes {
				if len(rt.Regions) > 0 && !slices.Contains(rt.Regions, region) {
					continue
				}

				name := fmt.Sprintf("AWS::%s::%s", svc.Name, rt.Name)
				input := cloudcontrol.ListResourcesInput{
					TypeName: &name,
				}
				paginator := cloudcontrol.NewListResourcesPaginator(ccClient, &input)
				for paginator.HasMorePages() {
					output, err := paginator.NextPage(ctx)
					if err != nil {
						log.Fatalf("failed to list resources for '%s': %v", rt.Name, err)
					}
					for _, res := range output.ResourceDescriptions {
						fmt.Printf("%s,%s,%s,%s\n", region, svc.Name, rt.Name, *res.Identifier)
					}
				}
			}
		}
	}
}

func getEnabledRegions(awsConfig aws.Config) []string {
	accountClient := account.NewFromConfig(
		awsConfig,
		func(o *account.Options) { o.Region = "us-east-1" },
	)

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
