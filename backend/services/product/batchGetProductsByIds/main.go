package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/RCGuo/aws-microservices-go/datasources/dynamodb/dbclient"
	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
)

func batchGetProducts(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("batchGetProducts")

  productKeys := models.ProductKeys{}
  err := json.Unmarshal([]byte(event.Body), &productKeys)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return api_utils.APIServerError(err)
  }
 
  var productsNumber int = len(productKeys)
  var reqList []map[string]types.AttributeValue
  var productList models.Products
  productTable := os.Getenv("DYNAMODB_PRODUCT_TABLE_NAME")
  for i, productKey := range productKeys {
    if err != nil {
      log.Printf("DynamoDB attributevalue MarshalMap error: %s\n", err.Error())
      return api_utils.APIServerError(err)
    }

    reqList = append(reqList, map[string]types.AttributeValue{
      "productId": &types.AttributeValueMemberS{Value: productKey.ProductId},
    },)

    if(len(reqList) == 100 || i+1 == productsNumber) {
      params := &dynamodb.BatchGetItemInput{
      	RequestItems: map[string]types.KeysAndAttributes{
          productTable: {
            Keys: reqList,
          },
        },
      }
    
      result, err := dbclient.Client.BatchGetItem(
        context.TODO(),
        params,
      )
      if err != nil {
        log.Printf("Create API call failed: %s", err)
        return api_utils.APIServerError(err)
      }

      products := models.Products{}
      err = attributevalue.UnmarshalListOfMaps(result.Responses[productTable], &products)
      if err != nil {
        log.Printf("Failed to unmarshal Dynamodb item, %s", err.Error())
        return api_utils.APIServerError(err)
      }

      productList = append(productList, products...)
      reqList = []map[string]types.AttributeValue{}
    }
  }

  return api_utils.APISuccessResponse(productList)
}

func main() {
  lambda.Start(batchGetProducts)
}
