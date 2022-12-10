package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/RCGuo/aws-microservices-go/datasources/dynamodb/dbclient"
	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
)

func updateProduct(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("updateProduct")

  keyList, gkErr := getKeyList(event.Body)
  if gkErr != nil {
    log.Printf("Error while getting update expression: %s\n", gkErr)
    return api_utils.APIServerError(gkErr)
  }

  productUpdate := models.ProductUpdate{}
  err := json.Unmarshal([]byte(event.Body), &productUpdate)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  updateData, err := attributevalue.MarshalMap(productUpdate)
	if err != nil {
		log.Println(err.Error())
		return api_utils.APIServerError(err)
	}

	key, err := attributevalue.MarshalMap(models.ProductKey{
    ProductId: event.PathParameters["productId"],
  })

	if err != nil {
		log.Println(err.Error())
		return api_utils.APIServerError(err)
	}

  params := &dynamodb.UpdateItemInput{
    TableName:                 aws.String(os.Getenv("DYNAMODB_PRODUCT_TABLE_NAME")),
  	Key:                       key,
  	ExpressionAttributeNames:  getExpressionAttributeNames(keyList),
  	ExpressionAttributeValues: getExpressionAttributeValues(keyList, updateData),
    UpdateExpression:          aws.String( getUpdateExpression(keyList) ),
    ReturnValues:              "UPDATED_NEW",
  }

  updateResult, err := dbclient.Client.UpdateItem(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("Update API call failed: %s", err)
    return api_utils.APIServerError(err)
  }

  log.Println("updateResult: ", updateResult)

  return api_utils.APISuccessResponse(updateResult)
}

func getKeyList(body string) ([]string, error) {
  keyList := []string{}
  data := map[string]json.RawMessage{}
  err := json.Unmarshal([]byte(body), &data)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return keyList, err
  }
  
  for key := range data {
    keyList = append(keyList, key)
  }

  return keyList, nil
}

func getUpdateExpression(list []string) string {
  updateExpression := []string{}
  for _, value := range list {
    updateExpression = append(updateExpression, fmt.Sprintf(" #%s = :%s", value, value)) 
  }
  
  return (fmt.Sprintf("SET %s", strings.Join(updateExpression, ",")))
}

func getExpressionAttributeNames(list []string) map[string]string {
  expressionAttributeNames := map[string]string{}
  for _, value := range list {
    expressionAttributeNames[fmt.Sprintf("#%s", value)] = value
  }
  
  return expressionAttributeNames
}

func getExpressionAttributeValues(list []string, data map[string]types.AttributeValue) map[string]types.AttributeValue {
  expressionAttributeValues := map[string]types.AttributeValue{}
  for _, value := range list {
    key := fmt.Sprintf(":%s", value)
    expressionAttributeValues[key] = data[key]
  }
  
  return expressionAttributeValues
}

func main() {
  lambda.Start(updateProduct)
}