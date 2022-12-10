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
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/RCGuo/aws-microservices-go/datasources/dynamodb/dbclient"
	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/RCGuo/aws-microservices-go/utils/auth_utils"
)

func updateCart(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("updateCart")

  cartItemUpdate := models.CartItemUpdate{}
  err := json.Unmarshal([]byte(event.Body), &cartItemUpdate)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  claims, err := auth_utils.ParseCognitoClaims(event)
  if err != nil {
    log.Printf("Fetching authentication data failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  newExpirationTime := time.Now().Add(72 * time.Hour).Unix()
  userId := claims.Username
  update := expression.UpdateBuilder{}.Set(
    expression.Name("quantity"), 
    expression.Value(cartItemUpdate.Quantity),
  )
  update = update.Set(
    expression.Name("expirationTime"), 
    expression.Value(newExpirationTime),
  )
  expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Failed to create dynamodb exprssion: %s", err.Error())
		return api_utils.APIServerError(err)
  }

  params := &dynamodb.UpdateItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_CART_TABLE_NAME")),
    Key: map[string]types.AttributeValue{
      "userId": &types.AttributeValueMemberS{Value: userId}, 
      "productId": &types.AttributeValueMemberS{Value: cartItemUpdate.ProductId}, 
    },
  	ExpressionAttributeNames:  expr.Names(),
  	ExpressionAttributeValues: expr.Values(),
  	UpdateExpression:          expr.Update(),
    ReturnValues: "UPDATED_NEW",
  }

  _, err = dbclient.Client.UpdateItem(
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
  lambda.Start(updateCart)
}
