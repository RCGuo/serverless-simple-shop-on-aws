package main

import (
	"context"
	"encoding/json"
	"errors"
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

func handleFavorite(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

  claims, err := auth_utils.ParseCognitoClaims(event)
  if err != nil {
    log.Printf("Fetching authentication data failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  userId := claims.Username
  if event.HTTPMethod == "GET" {
    return getUserAllFavorites(userId)
  } else if event.HTTPMethod == "POST" {
    favorite := models.UpdateProductFavorite{}
    err := json.Unmarshal([]byte(event.Body), &favorite)
    if err != nil {
      log.Printf("JSON unmarshal failed: %s\n", err)
      return api_utils.APIServerError(err)
    }
  
    if favorite.Favorite {
      return addToFavorite(favorite, userId)
    } else if ! favorite.Favorite {
      return deleteFromFavorite(favorite, userId)
    } else {
      return api_utils.APIBadRequest(errors.New("unknown request"))
    }
  } else {
    return api_utils.APIBadRequest(errors.New("unknown request method"))
  }
}

func getUserAllFavorites(userId string) (events.APIGatewayProxyResponse, error) {
	log.Println("getUserAllFavorites")

  keyCond := expression.Key("userId").Equal(expression.Value(userId))
  expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		log.Printf("Failed to create dynamodb exprssion: %s", err.Error())
		return api_utils.APIServerError(err)
  }

	params := &dynamodb.QueryInput{
		TableName: aws.String(os.Getenv("DYNAMODB_FAVORITE_TABLE_NAME")),
    ExpressionAttributeNames:  expr.Names(),
    ExpressionAttributeValues: expr.Values(),
    KeyConditionExpression:    expr.KeyCondition(),
	}

	data, err := dbclient.Client.Query(
		context.TODO(),
		params,
	)
	if err != nil {
		log.Printf("DynamoDB Scan API call failed: %s", err)
		return api_utils.APIServerError(err)
	}

  favorites := models.ProductFavorites{}
	err = attributevalue.UnmarshalListOfMaps(data.Items, &favorites)
	if err != nil {
		log.Printf("Failed to unmarshal Dynamodb items: %s", err.Error())
		return api_utils.APIServerError(err)
  }
  
	return api_utils.APISuccessResponse(favorites)
}

func addToFavorite(favorite models.UpdateProductFavorite, userId string) (events.APIGatewayProxyResponse, error) {
  log.Println("addFavorite")

  favorite.UserId = userId
  params := &dynamodb.PutItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_FAVORITE_TABLE_NAME")),
    // Item: data,
    Item: map[string]types.AttributeValue{
      "userId": &types.AttributeValueMemberS{Value: favorite.UserId},
      "productId": &types.AttributeValueMemberS{Value: favorite.ProductId},
    },
  }

  _, err := dbclient.Client.PutItem(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("Create API call failed: %s", err)
    return api_utils.APIServerError(err)
  }

  return api_utils.APISuccessResponse(nil)
}

func deleteFromFavorite(favorite models.UpdateProductFavorite, userId string) (events.APIGatewayProxyResponse, error) {
  log.Println("deleteFromFavorite")

  params := &dynamodb.DeleteItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_FAVORITE_TABLE_NAME")),
    Key: map[string]types.AttributeValue{
      "userId": &types.AttributeValueMemberS{Value: userId}, 
      "productId": &types.AttributeValueMemberS{Value: favorite.ProductId}, 
    },
  }

  deleteResult, err := dbclient.Client.DeleteItem(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("DynamoDB delete API call failed: %s", err)
    return api_utils.APIServerError(err)
  }

  log.Println("DynamoDB delete result: ", deleteResult)

  return api_utils.APISuccessResponse(nil)
}

func main() {
  lambda.Start(handleFavorite)
}