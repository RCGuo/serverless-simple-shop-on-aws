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

	"github.com/RCGuo/aws-microservices-go/datasources/dynamodb/dbclient"
	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/RCGuo/aws-microservices-go/utils/auth_utils"
)

func getPastOrders(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("getOrdersByUserId")

  claims, err := auth_utils.ParseCognitoClaims(event)
  if err != nil {
    log.Printf("Fetching authentication data failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  userId := claims.Username
  keyCond := expression.Key("userId").Equal(expression.Value(userId))
  expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		log.Printf("Failed to create dynamodb exprssion: %s", err.Error())
		return api_utils.APIServerError(err)
  }

  params := &dynamodb.QueryInput{
  	TableName:                 aws.String(os.Getenv("DYNAMODB_ORDER_TABLE_NAME")),
  	ExpressionAttributeNames:  expr.Names(),
  	ExpressionAttributeValues: expr.Values(),
  	KeyConditionExpression:    expr.KeyCondition(),
  }

  data, err := dbclient.Client.Query(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("DynamoDB Query API call failed: %s", err)
    return api_utils.APIServerError(err)
  }

  orders := models.Orders{}
  err = attributevalue.UnmarshalListOfMaps(data.Items, &orders)
  if err != nil {
    log.Printf("Failed to unmarshal Dynamodb item, %s", err.Error())
    return api_utils.APIServerError(err)
  }

  return api_utils.APISuccessResponse(orders)
}

func main(){
  lambda.Start(getPastOrders)
}