package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/RCGuo/aws-microservices-go/datasources/opensearch/opclient"
	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/dynamodb_utils"
)

const indexName = "product-sold-counter"
func streamOrderToOpenSearch(ctx context.Context, event events.DynamoDBEvent) {
  log.Println("streamOrderToOpenSearch")

  var body strings.Builder
  for _, record := range event.Records {
		item, err := dynamodb_utils.FromDynamoDBEventAVMap(record.Change.NewImage)
		if err != nil {
      log.Printf("FromDynamoDBEventAVMap failed: %s\n", err.Error())
		}

    order := models.Order{}
    err = attributevalue.UnmarshalMap(item, &order)
    if err != nil {
      log.Printf("UnmarshalMap failed: %s\n", err.Error())
    }

    if order.PaymentStatus != "charge.succeeded" {
      continue
    }
    log.Printf("# order: %+v", order)
    for _, item := range order.Items {
      if err != nil {
        log.Printf("json marshel failed: %s\n", err.Error())
      }
  
      if item.Quantity <= 0 {
        log.Printf("item quantity cannot be less than or equal to zero")
      }

      body.WriteString(`{"update" : { "_index" : "` + indexName + `", "_id" : "` + item.ProductId + `" }}`)
      body.WriteString("\n")
      body.WriteString(
        `{` +
          `"script": {` +
            `"source": "ctx._source.counter += params.count",` + 
            `"lang": "painless",` +
            `"params": {` +
              `"count":` + strconv.Itoa(item.Quantity) +
            `}` +
          `},` +
          `"upsert": {` +
            `"counter":` + strconv.Itoa(item.Quantity) +
          `}` +
        `}`)
      body.WriteString("\n")

      log.Println("# body string: ", body.String())

      _, err = opclient.Client.Bulk(
        strings.NewReader(body.String()),
      )
      if err != nil {
        log.Println("failed to perform bulk operations", err)
      }
    }
  }
}

func main() {
	lambda.Start(streamOrderToOpenSearch)
}