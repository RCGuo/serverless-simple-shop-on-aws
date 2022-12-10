package main

import (
	"context"
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
)

func productHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Printf("# event : %+v\n", event)

  if event.QueryStringParameters != nil {
    if event.QueryStringParameters["category"] != "" {
      category := event.QueryStringParameters["category"]
      return getProductByCategory(category)      
    } else if event.QueryStringParameters["topic"] != "" {
      topic := event.QueryStringParameters["topic"]
      return getProductByTopic(topic) 
    } else if event.QueryStringParameters["productId"] != "" {
      productId := event.QueryStringParameters["productId"]
      return getProductById(productId)
    } else {
      return api_utils.APIBadRequest(errors.New("query parameter of product is not valid"))
    }
  } else {
    return getAllProducts()
  }
}

func getProductById(productId string) (events.APIGatewayProxyResponse, error) {
  log.Println("getProductById")

  params := &dynamodb.GetItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_PRODUCT_TABLE_NAME")),
    Key: map[string]types.AttributeValue{
      "productId": &types.AttributeValueMemberS{Value: productId}, 
    },
  }

  data, err := dbclient.Client.GetItem(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("DynamoDB GetItem API call failed: %s", err)
    return api_utils.APIErrResponse(400, err)
  }

  product := models.Product{}
  err = attributevalue.UnmarshalMap(data.Item, &product)
  if err != nil {
    log.Printf("Failed to unmarshal Dynamodb item, %s", err.Error())
    return api_utils.APIServerError(err)
  }

  return api_utils.APISuccessResponse(product)
}

func getAllProducts() (events.APIGatewayProxyResponse, error) {
	log.Println("getAllProduct")

	params := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DYNAMODB_PRODUCT_TABLE_NAME")),
	}

	data, err := dbclient.Client.Scan(
		context.TODO(),
		params,
	)
	if err != nil {
		log.Printf("DynamoDB Scan API call failed: %s", err)
		return api_utils.APIServerError(err)
	}

  products := []models.Product{}
	err = attributevalue.UnmarshalListOfMaps(data.Items, &products)
	if err != nil {
		log.Printf("Failed to unmarshal Dynamodb items: %s", err.Error())
		return api_utils.APIServerError(err)
  }
  
	return api_utils.APISuccessResponse(products)
}

func getProductByCategory(category string) (events.APIGatewayProxyResponse, error) {
  log.Println("getProductByCategory")

  keyCond := expression.Key("category").Equal(expression.Value(category))
  expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		log.Printf("Failed to create dynamodb exprssion: %s", err.Error())
		return api_utils.APIServerError(err)
  }

  params := &dynamodb.QueryInput{
    TableName:                 aws.String(os.Getenv("DYNAMODB_PRODUCT_TABLE_NAME")),
    ExpressionAttributeNames:  expr.Names(),
    ExpressionAttributeValues: expr.Values(),
    KeyConditionExpression:    expr.KeyCondition(),
    IndexName:                 aws.String("category-index"),
  }

  queryResult, err := dbclient.Client.Query(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("Failed to query table, %s", err.Error())
    return api_utils.APIServerError(err)
  }

  products := models.Products{}
  err = attributevalue.UnmarshalListOfMaps(queryResult.Items, &products)
  if err != nil {
    log.Printf("Failed to unmarshal Dynamodb item [getProduct], %s", err.Error())
    return api_utils.APIServerError(err)
  }

  return api_utils.APISuccessResponse(products)
}

func getProductByTopic(topic string) (events.APIGatewayProxyResponse, error) {
  log.Println("getProductByCategory")

  keyCond := expression.Key("topic").Equal(expression.Value(topic))
  expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		log.Printf("Failed to create dynamodb exprssion: %s", err.Error())
		return api_utils.APIServerError(err)
  }

  params := &dynamodb.QueryInput{
    TableName:                 aws.String(os.Getenv("DYNAMODB_PRODUCT_TABLE_NAME")),
    ExpressionAttributeNames:  expr.Names(),
    ExpressionAttributeValues: expr.Values(),
    KeyConditionExpression:    expr.KeyCondition(),
    IndexName:                 aws.String("topic-index"),
  }

  queryResult, err := dbclient.Client.Query(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("Failed to query table, %s", err.Error())
    return api_utils.APIServerError(err)
  }

  products := models.Products{}
  err = attributevalue.UnmarshalListOfMaps(queryResult.Items, &products)
  if err != nil {
    log.Printf("Failed to unmarshal Dynamodb item [getProduct], %s", err.Error())
    return api_utils.APIServerError(err)
  }

  return api_utils.APISuccessResponse(products)
}

func main() {
	lambda.Start(productHandler)
}