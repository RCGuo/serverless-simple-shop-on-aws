package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/opensearch-project/opensearch-go/opensearchapi"

	"github.com/RCGuo/aws-microservices-go/datasources/opensearch/opclient"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/RCGuo/aws-microservices-go/utils/opensearch_utils"
)

func fuzzyQuery(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  var query, index, sizeStr string
  if event.QueryStringParameters != nil {
    if event.QueryStringParameters["query"] != "" {
      query = event.QueryStringParameters["query"]
    }
    if event.QueryStringParameters["index"] != "" {
      index = event.QueryStringParameters["index"]
    }
    if event.QueryStringParameters["size"] != "" {
      sizeStr = event.QueryStringParameters["size"]
    } else {
      sizeStr = "100"
    }
  }
  
  size, err := strconv.Atoi(sizeStr)
  if err != nil {
    log.Println("Invalid type of size parameter: ", err.Error())
    return api_utils.APIBadRequest(err)
  }

  content := strings.NewReader(fmt.Sprintf(`{
    "query": {
      "simple_query_string": {
        "query": "%s~",
        "fields": ["name"],
        "flags": "ALL",
        "fuzzy_transpositions": true,
        "fuzzy_max_expansions": 50,
        "fuzzy_prefix_length": 0,
        "minimum_should_match": 1,
        "default_operator": "or",
        "analyzer": "standard",
        "lenient": false,
        "quote_field_suffix": "",
        "analyze_wildcard": false,
        "auto_generate_synonyms_phrase_query": true
      }
    }
  }`, query ))

  search := opensearchapi.SearchRequest{
    Index: []string{index},
    Body: content,
    Size: &size,
    Pretty: true,
    Human: true,
  }

  response, err := search.Do(context.Background(), opclient.Client)
  if err != nil {
    log.Println("failed to search document ", err.Error())
    return api_utils.APIServerError(err)
  }

  log.Printf("searchResponse: %+v", response)

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
	lambda.Start(fuzzyQuery)
}