package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
)

func Scan(ctx context.Context, cfg *Config) ([]Resource, error) {
	awsConfig, err := config.LoadDefaultConfig(
		ctx,
		config.WithAppID("github.com/mkuzmin/dritf"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config: %w", err)
	}

	var resources []Resource
	for _, region := range cfg.Regions {
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
						return nil, fmt.Errorf("failed to list resources for '%s' (%s): %v", name, region, err)
					}

					for _, res := range output.ResourceDescriptions {
						id := *res.Identifier
						resources = append(resources, Resource{
							Region:   region,
							Service:  service.Name,
							TypeName: resourceType.Name,
							Id:       id,
						})
					}
				}
			}
		}
	}

	return resources, nil
}
