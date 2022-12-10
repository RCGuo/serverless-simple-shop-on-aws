package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/opensearch-project/opensearch-go/opensearchapi"

	"github.com/RCGuo/aws-microservices-go/datasources/opensearch/opclient"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/RCGuo/aws-microservices-go/utils/opensearch_utils"
)

func purchasedRanking(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  var index, sizeStr, sortField, sortDirect string
  if event.QueryStringParameters != nil {
    if event.QueryStringParameters["index"] != "" {
      index = event.QueryStringParameters["index"]
    }
    if event.QueryStringParameters["size"] != "" {
      sizeStr = event.QueryStringParameters["size"]
    } else {
      sizeStr = "100"
    }
    if event.QueryStringParameters["sortField"] != "" {
      sortField = event.QueryStringParameters["sortField"]
    }
    if event.QueryStringParameters["sortDirect"] != "" {
      sortDirect = event.QueryStringParameters["sortDirect"]
    } else {
      sortDirect = "desc"
    }
  }

  size, err := strconv.Atoi(sizeStr)
  if err != nil {
    log.Println("Invalid type of size parameter: ", err.Error())
    return api_utils.APIBadRequest(err)
  }

  if index == "" || sortField == ""  {
    log.Println("Invalid index or field name")
    return api_utils.APIBadRequest(fmt.Errorf("%s", "invalid index or field name"))
  }

  search := opensearchapi.SearchRequest{
    Index: []string{index},
    Size: &size,
    Sort: []string{fmt.Sprintf("%s:%s", sortField, sortDirect)},
  }

  response, err := search.Do(context.Background(), opclient.Client)
  if err != nil {
    log.Println("failed to search document ", err.Error())
    return api_utils.APIServerError(err)
  }

  if response.IsError() {
    res, err := opensearch_utils.ParseErrorResponse(response)
    if err != nil {
      log.Println("failed to parse response ", err.Error())
      return api_utils.APIServerError(err)
    }
    return api_utils.APIResponse(response.StatusCode, res)
  } else {
    res, err := opensearch_utils.ParseResponse(response)
    if err != nil {
      log.Println("failed to parse response ", err.Error())
      return api_utils.APIServerError(err)
    }
    return api_utils.APIResponse(response.StatusCode, res)
  }

}

func main() {
  lambda.Start(purchasedRanking)
}