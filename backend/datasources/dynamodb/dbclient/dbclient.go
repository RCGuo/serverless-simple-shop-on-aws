package dbclient

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var Client *dynamodb.Client

func init(){
  optFuns := []func(*config.LoadOptions) error {
    config.WithRegion("us-east-1"),
  }

  dynamoDBLocal := os.Getenv("DYNAMODB_LOCAL")
  if strings.HasPrefix(dynamoDBLocal, "http") {
    customResolver := aws.EndpointResolverWithOptionsFunc(
      func(service, region string, options ...interface{}) (aws.Endpoint, error) {
        return aws.Endpoint{
					URL: dynamoDBLocal,
        }, nil
    })
    optFuns = append(optFuns, config.WithEndpointResolverWithOptions(customResolver))
  } 

  cfg, err := config.LoadDefaultConfig(context.TODO(), optFuns...)
  if err != nil {
    log.Printf("Error initializing client with DynamoDB: %s\n", err.Error())
  }
  Client = dynamodb.NewFromConfig(cfg)
}
