package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/aws/aws-sdk-go-v2/service/account/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	ccTypes "github.com/aws/aws-sdk-go-v2/service/cloudcontrol/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"log"
	"slices"
)

func scanTypes(cfg *Config) {
	awsConfig := setupAWSConfig()
	accountId := getAccountId(awsConfig)
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
				resources := scanResourceType(ccClient, svc, rt)

				for _, res := range resources {
					fmt.Printf("%s,%s,%s,%s,%s\n", accountId, region, svc.Name, res.TypeName, res.Id)
				}
			}
		}
	}
}

func setupAWSConfig() aws.Config {
	awsConfig, err := config.LoadDefaultConfig(
		ctx,
		config.WithAppID("github.com/mkuzmin/dritf"),
	)
	if err != nil {
		log.Fatal("failed to create AWS client: ", err)
	}
	return awsConfig
}

func getAccountId(awsConfig aws.Config) string {
	stsClient := sts.NewFromConfig(
		awsConfig,
		func(o *sts.Options) { o.Region = "us-east-1" },
	)

	result, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalf("unable to get caller identity: %v", err)
	}

	return *result.Account
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

type Resource struct {
	TypeName string
	Id       string
}

func scanResourceType(ccClient *cloudcontrol.Client, svc Service, rt ResourceType) []Resource {
	name := fmt.Sprintf("AWS::%s::%s", svc.Name, rt.Name)
	input := cloudcontrol.ListResourcesInput{TypeName: &name}

	resources := listResources(ccClient, &input, rt.Name)
	var result []Resource
	for _, res := range resources {
		result = append(result, Resource{
			TypeName: rt.Name,
			Id:       *res.Identifier,
		})

		for _, depType := range rt.DependentTypes {
			depResources := listDependentResources(ccClient, &res, svc, depType)
			result = append(result, depResources...)
		}
	}
	return result
}

func listResources(
	ccClient *cloudcontrol.Client,
	input *cloudcontrol.ListResourcesInput,
	resourceTypeName string,
) []ccTypes.ResourceDescription {
	var resources []ccTypes.ResourceDescription
	paginator := cloudcontrol.NewListResourcesPaginator(ccClient, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to list resources for '%s': %v", resourceTypeName, err)
		}
		resources = append(resources, output.ResourceDescriptions...)
	}
	return resources
}

func listDependentResources(ccClient *cloudcontrol.Client, res *ccTypes.ResourceDescription, svc Service, depType DepType) []Resource {
	name := fmt.Sprintf("AWS::%s::%s", svc.Name, depType.Name)
	var model string
	if depType.Property == nil {
		model = fmt.Sprintf(
			`{"%s": "%s"}`,
			depType.Ref,
			*res.Identifier,
		)
	} else {
		model = fmt.Sprintf(
			`{"%s": "%s"}`,
			depType.Ref,
			getJsonProperty(*res.Properties, *depType.Property),
		)
	}
	input := cloudcontrol.ListResourcesInput{
		TypeName:      &name,
		ResourceModel: &model,
	}
	resources := listResources(ccClient, &input, depType.Name)
	var result []Resource
	for _, res := range resources {
		result = append(result, Resource{
			TypeName: depType.Name,
			Id:       *res.Identifier,
		})
	}
	return result
}

func getJsonProperty(doc string, property string) string {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(doc), &result)
	if err != nil {
		log.Fatal("failed to unmarshal json: ", err)
	}

	return result[property].(string)
}
