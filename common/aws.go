package common

import (
	"context"
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

func ScanTypes(ctx context.Context, cfg *Config) chan Resource {
	resourceChan := make(chan Resource)
	go scanTypesInternal(ctx, cfg, resourceChan)
	return resourceChan
}

func scanTypesInternal(ctx context.Context, cfg *Config, resourceChan chan Resource) {

	awsConfig := setupAWSConfig(ctx)
	accountId := getAccountId(ctx, awsConfig)
	enabledRegions := getEnabledRegions(ctx, awsConfig)

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
				scanResourceType(ctx, ccClient, accountId, region, svc, rt, resourceChan)
			}
		}
	}
	close(resourceChan)
}

func setupAWSConfig(ctx context.Context) aws.Config {
	awsConfig, err := config.LoadDefaultConfig(
		ctx,
		config.WithAppID("github.com/mkuzmin/dritf"),
	)
	if err != nil {
		log.Fatal("failed to create AWS client: ", err)
	}
	return awsConfig
}

func getAccountId(ctx context.Context, awsConfig aws.Config) string {
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

func getEnabledRegions(ctx context.Context, awsConfig aws.Config) []string {
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
	Account  string
	Region   string
	Service  string
	TypeName string
	Id       string
}

func scanResourceType(ctx context.Context, ccClient *cloudcontrol.Client, accountId string, region string, svc Service, rt ResourceType, outChan chan Resource) {
	name := fmt.Sprintf("AWS::%s::%s", svc.Name, rt.Name)
	input := cloudcontrol.ListResourcesInput{TypeName: &name}

	resources := listResources(ctx, ccClient, &input, svc, rt.Name, accountId, region, outChan)
	for _, res := range resources {
		for _, depType := range rt.DependentTypes {
			listDependentResources(ctx, ccClient, &res, svc, depType, accountId, region, outChan)
		}
	}
}

func listResources(ctx context.Context, ccClient *cloudcontrol.Client, input *cloudcontrol.ListResourcesInput, svc Service, typeName string, accountId string, region string, outChan chan Resource) []ccTypes.ResourceDescription {
	var resources []ccTypes.ResourceDescription
	paginator := cloudcontrol.NewListResourcesPaginator(ccClient, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to list resources for '%s': %v", typeName, err)
		}
		for _, res := range output.ResourceDescriptions {
			outChan <- Resource{
				Account:  accountId,
				Region:   region,
				Service:  svc.Name,
				TypeName: typeName,
				Id:       *res.Identifier,
			}
			resources = append(resources, res)
		}
	}
	return resources
}

func listDependentResources(ctx context.Context, ccClient *cloudcontrol.Client, res *ccTypes.ResourceDescription, svc Service, depType DepType, accountId string, region string, outChan chan Resource) {
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

	listResources(ctx, ccClient, &input, svc, depType.Name, accountId, region, outChan)
}

func getJsonProperty(doc string, property string) string {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(doc), &result)
	if err != nil {
		log.Fatal("failed to unmarshal json: ", err)
	}

	return result[property].(string)
}
