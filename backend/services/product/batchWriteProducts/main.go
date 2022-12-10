package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	"github.com/RCGuo/aws-microservices-go/utils/auth_utils"
)

func batchWriteProducts(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("batchWriteProduct")

  isAdmin, err := auth_utils.IsRequestByCognitoAdmin(event)
  if err != nil {
    log.Println(err.Error())
    return api_utils.APIServerError(err)
  }

  if !isAdmin {
    log.Printf("insufficient permissions: %s", err)
    return api_utils.APIServerError(fmt.Errorf("insufficient permissions: %s", err))
  }

  products := models.Products{}
  err = json.Unmarshal([]byte(event.Body), &products)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return api_utils.APIServerError(err)
  }
 
  var productLength int = len(products)
  var writeReqs []types.WriteRequest
  for i, product := range products {
    product.ImageFile = fmt.Sprintf("%s/%s/%s", os.Getenv("IMAGE_CDN"), product.Category, product.ImageFile)
    product.ProductId = base64.StdEncoding.EncodeToString([]byte(product.Name))
    item, err := attributevalue.MarshalMap(product)

    if err != nil {
      log.Printf("DynamoDB attributevalue MarshalMap error: %s\n", err.Error())
      return api_utils.APIServerError(err)
    }

    writeReqs = append(writeReqs, types.WriteRequest{
      PutRequest: &types.PutRequest{Item: item},
    })

    if(len(writeReqs) == 25 || i+1 == productLength) {
      params := &dynamodb.BatchWriteItemInput{
        RequestItems: map[string][]types.WriteRequest{
          os.Getenv("DYNAMODB_PRODUCT_TABLE_NAME"): writeReqs,
        },
      }
    
      batchWriteResult, err := dbclient.Client.BatchWriteItem(
        context.TODO(),
        params,
      )
      log.Println("batchWriteResult: ", batchWriteResult)
  
      if err != nil {
        log.Printf("Create API call failed: %s", err)
        return api_utils.APIServerError(err)
      }
      writeReqs = []types.WriteRequest{}
    }
  }

  return api_utils.APISuccessResponse(nil)
}

func main() {
  lambda.Start(batchWriteProducts)
}
