package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/RCGuo/aws-microservices-go/datasources/dynamodb/dbclient"
	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/RCGuo/aws-microservices-go/utils/auth_utils"
)

func cartHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Printf("# event : %+v\n", event)
  
  claims, err := auth_utils.ParseCognitoClaims(event)
  if err != nil {
    log.Printf("Fetching authentication data failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  userId := claims.Username
  if event.PathParameters != nil {
    productId := event.PathParameters["productId"]
    log.Println("# userId: ", userId)
    return getCartItemByProductId(userId, productId)
  } else {
    return getCartItems(userId)
  }
}

func getCartItems(userId string) (events.APIGatewayProxyResponse, error) {
  log.Println("getUserCartItems")

  keyCond := expression.Key("userId").Equal(expression.Value(userId))
  expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		log.Printf("Failed to create dynamodb exprssion: %s", err.Error())
		return api_utils.APIServerError(err)
  }

  params := &dynamodb.QueryInput{
    TableName: aws.String(os.Getenv("DYNAMODB_CART_TABLE_NAME")),
    ExpressionAttributeNames:  expr.Names(),
    ExpressionAttributeValues: expr.Values(),
    KeyConditionExpression:    expr.KeyCondition(),
  }
  
  data, err := dbclient.Client.Query(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("DynamoDB Query API call failed: %s", err.Error())
    return api_utils.APIServerError(err)
  }

  carts := models.CartItems{}
  err = attributevalue.UnmarshalListOfMaps(data.Items, &carts)
  if err != nil {
    log.Printf("Failed to unmarshal Dynamodb item, %s", err.Error())
    return api_utils.APIServerError(err)
  }

  return api_utils.APISuccessResponse(carts)
}

func getCartItemByProductId(userId, productId string) (events.APIGatewayProxyResponse, error) {
  log.Println("getUserCartItemByProductId")

  params := &dynamodb.GetItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_CART_TABLE_NAME")),
    Key: map[string]types.AttributeValue{
      "userId": &types.AttributeValueMemberS{Value: userId}, 
      "productId": &types.AttributeValueMemberS{Value: productId}, 
    },
  }
  
  data, err := dbclient.Client.GetItem(
    context.TODO(),
    params,
  )

  if err != nil {
    log.Printf("DynamoDB Query API call failed: %s", err.Error())
    return api_utils.APIServerError(err)
  }

  cart := models.CartItem{}
  err = attributevalue.UnmarshalMap(data.Item, &cart)
  if err != nil {
    log.Printf("Failed to unmarshal Dynamodb item, %s", err.Error())
    return api_utils.APIServerError(err)
  }

  return api_utils.APISuccessResponse(cart)
}

func main() {
	lambda.Start(cartHandler)
}