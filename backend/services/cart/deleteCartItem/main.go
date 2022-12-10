package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/RCGuo/aws-microservices-go/datasources/dynamodb/dbclient"
	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/RCGuo/aws-microservices-go/utils/auth_utils"
)

func deleteCartItemById(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("deleteProduct")

  cartItemDelete := models.CartItemDelete{}
  err := json.Unmarshal([]byte(event.Body), &cartItemDelete)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  claims, err := auth_utils.ParseCognitoClaims(event)
  if err != nil {
    log.Printf("Fetching authentication data failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  userId := claims.Username
  productId := cartItemDelete.ProductId
  params := &dynamodb.DeleteItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_CART_TABLE_NAME")),
    Key: map[string]types.AttributeValue{
      "userId": &types.AttributeValueMemberS{Value: userId}, 
      "productId": &types.AttributeValueMemberS{Value: productId}, 
    },
  }

  _, err = dbclient.Client.DeleteItem(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("DynamoDB delete API call failed: %s", err)
    return api_utils.APIServerError(err)
  }

  return api_utils.APISuccessResponse(nil)
}


func main() {
  lambda.Start(deleteCartItemById)
}