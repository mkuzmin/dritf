package main

import (
	"bufio"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var ctx = context.Background()

func main() {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithAppID("mkuzmin.io/dritf"))
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}
	cfg.Region = "eu-central-1" // TODO iterate over all regions
	cc := cloudcontrol.NewFromConfig(cfg)

	f, err := os.Open("types.txt")
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	//for {
	//	s, _ := r.ReadString('\n')
	//	if s == "---\n" {
	//		break
	//	}
	//}

	for {
		s, err := r.ReadString('\n')
		if err == io.EOF {
			break // TODO bug: last line may contain data
		}
		if err != nil {
			log.Fatal("Error reading file: ", err)
		}
		s = strings.TrimSuffix(s, "\n")
		if strings.HasPrefix(s, "#") {
			continue
		}

		print("=== " + s + "\r")

		input := cloudcontrol.ListResourcesInput{
			TypeName: &s,
		}
		paginator := cloudcontrol.NewListResourcesPaginator(cc, &input)

		for paginator.HasMorePages() {
			time.Sleep(200 * time.Millisecond) // prevent API throttling
			output, err := paginator.NextPage(ctx)
			if err != nil {
				log.Fatal("Failed to list types: ", err)
			}

			if len(output.ResourceDescriptions) > 0 {
				println()
				for _, res := range output.ResourceDescriptions {
					println(*res.Identifier)
				}
			}
		}
	}
}
