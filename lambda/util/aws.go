package util

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadDefaultConfig(ctx context.Context, region, endpoint, service string) (cfg aws.Config, err error) {
	return config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
		o.Region = region
		if endpoint != "" {
			o.EndpointResolverWithOptions = createEndpointResolverWithOptions(service, region, endpoint)
		}
		return nil
	})
}

func createEndpointResolverWithOptions(service, region, endpoint string) aws.EndpointResolverWithOptions {
	return aws.EndpointResolverWithOptionsFunc(func(s, r string, opts ...interface{}) (aws.Endpoint, error) {
		if s == service && r == region {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})
}
