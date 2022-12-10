package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/RCGuo/aws-microservices-go/datasources/opensearch/opclient"
	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/RCGuo/aws-microservices-go/utils/auth_utils"
)

const indexName = "product"
func bullWriteProductIndexes(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("handleDynamoDBStream")

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

  var body strings.Builder
  for _, product := range products {
    product.ImageFile = fmt.Sprintf("%s/%s/%s", os.Getenv("IMAGE_CDN"), product.Category, product.ImageFile)
    product.ProductId = base64.StdEncoding.EncodeToString([]byte(product.Name))

    marshelJson, err := json.Marshal(product)
    if err != nil {
      log.Printf("json marshel failed: %s\n", err.Error())
      return api_utils.APIServerError(err)
    }

    body.WriteString(`{"create" : { "_index" : "` + indexName + `", "_id" : "` + product.ProductId + `" }}`)
    body.WriteString("\n")
    body.Write(marshelJson)
    body.WriteString("\n")

    log.Println(body.String())
  }

  log.Println(body.String())
  
  _, err = opclient.Client.Bulk(
		strings.NewReader(body.String()),
	)
  if err != nil {
    log.Println("failed to perform bulk operations", err)
  }

  return api_utils.APISuccessResponse(nil)
}

func main() {
	lambda.Start(bullWriteProductIndexes)
}