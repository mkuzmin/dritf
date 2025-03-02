package terraform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type Resources []Resource

type Resource struct {
	Type    string
	Address string
	Id      string
}

func GetResources(ctx context.Context, workDir string) (*Resources, error) {
	tf, err := tfexec.NewTerraform(workDir, "terraform")
	if err != nil {
		return nil, fmt.Errorf("error running NewTerraform: %s", err)
	}

	show, err := tf.Show(ctx)
	if err != nil {
		return nil, fmt.Errorf("error running Show: %s", err)
	}

	if show.Values == nil {
		return nil, fmt.Errorf("there is no resources in the state")
	}

	resources := collect(show.Values.RootModule)
	return &resources, nil
}

func collect(module *tfjson.StateModule) Resources {
	var resources Resources

	for _, r := range module.Resources {
		if r.Mode != "managed" {
			continue
		}

		if r.ProviderName != "registry.terraform.io/hashicorp/aws" {
			continue
		}

		resources = append(resources, Resource{
			Type:    r.Type,
			Address: r.Address,
			Id:      r.AttributeValues["id"].(string),
		})
	}

	for _, m := range module.ChildModules {
		resources = append(resources, collect(m)...)
	}

	return resources
}
