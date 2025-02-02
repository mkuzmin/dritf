package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
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

	for _, region := range cfg.Regions {
		println("=== ", region)
		ccClient := cloudcontrol.NewFromConfig(
			awsConfig,
			func(o *cloudcontrol.Options) { o.Region = region },
		)

		for _, svc := range cfg.Services {
			println("   === ", svc.Name)

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
						println(*res.Identifier)
					}
				}
			}
			println()
		}
		println()
	}
}
