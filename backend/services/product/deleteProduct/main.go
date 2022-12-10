package main

import (
	"context"
	"log"
	"os"

	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/RCGuo/aws-microservices-go/datasources/dynamodb/dbclient"
)

func deleteProduct(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("deleteProduct")

  productId := event.PathParameters["productId"]
  params := &dynamodb.DeleteItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_PRODUCT_TABLE_NAME")),
    Key: map[string]types.AttributeValue{
      "productId": &types.AttributeValueMemberS{Value: productId}, 
    },
  }

  _, err := dbclient.Client.DeleteItem(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("Delete API call failed: %s", err)
    return api_utils.APIServerError(err)
  }

  return api_utils.APISuccessResponse(nil)
}


func main() {
  lambda.Start(deleteProduct)
}