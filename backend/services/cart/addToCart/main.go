package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/RCGuo/aws-microservices-go/datasources/dynamodb/dbclient"
	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/RCGuo/aws-microservices-go/utils/auth_utils"
)

func addToCart(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("addToCart")

  cartItem := models.CartItem{}
  err := json.Unmarshal([]byte(event.Body), &cartItem)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  claims, err := auth_utils.ParseCognitoClaims(event)
  if err != nil {
    log.Printf("Fetching authentication data failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  cartItem.UserId = claims.Username
  cartItem.ExpirationTime = time.Now().Add(72 * time.Hour).Unix()
	data, err := attributevalue.MarshalMap(cartItem)
	if err != nil {
		log.Printf("DynamoDB attributevalue MarshalMap error: %s\n", err.Error())
    return api_utils.APIServerError(err)
	}

  params := &dynamodb.PutItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_CART_TABLE_NAME")),
    Item: data,
  }

  _, err = dbclient.Client.PutItem(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("Create API call failed: %s", err)
    return api_utils.APIServerError(err)
  }

  return api_utils.APISuccessResponse(nil)
}

func main() {
  lambda.Start(addToCart)
}
