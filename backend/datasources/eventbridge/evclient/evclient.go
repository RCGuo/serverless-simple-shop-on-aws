package evclient

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
)
	
var Client *eventbridge.Client

func init(){
  optFuns := []func(*config.LoadOptions) error {
    config.WithRegion("us-east-1"),
  }

  cfg, err := config.LoadDefaultConfig(context.TODO(), optFuns...)
  if err != nil {
    log.Printf("Error initializing client with DynamoDB: %s\n", err.Error())
  }
  Client = eventbridge.NewFromConfig(cfg)
}