package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/RCGuo/aws-microservices-go/datasources/dynamodb/dbclient"
	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
)

func createProduct(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("createProduct")

  product := models.Product{}
  err := json.Unmarshal([]byte(event.Body), &product)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return api_utils.APIServerError(err)
  }
 
  product.ProductId = base64.StdEncoding.EncodeToString([]byte(product.Name))
	data, err := attributevalue.MarshalMap(product)
	if err != nil {
		log.Printf("DynamoDB attributevalue MarshalMap error: %s\n", err.Error())
    return api_utils.APIServerError(err)
	}
  
  condition := expression.Name("productId").NotEqual(expression.Value(product.ProductId))
  exp, _ := expression.NewBuilder().WithCondition(condition).Build()

  params := &dynamodb.PutItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_PRODUCT_TABLE_NAME")),
    Item: data,
    ConditionExpression: exp.Condition(),
    ExpressionAttributeValues: exp.Values(),
    ExpressionAttributeNames:  exp.Names(),
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
  lambda.Start(createProduct)
}
