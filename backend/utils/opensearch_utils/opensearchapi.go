package opensearch_utils

import (
	"encoding/json"
	"log"

	"github.com/opensearch-project/opensearch-go/opensearchapi"

	"github.com/RCGuo/aws-microservices-go/models"
)

func ParseResponse(response *opensearchapi.Response) (*[]models.SearchResponseHit, error) {

  var res models.SearchResponse
  err := json.NewDecoder(response.Body).Decode(&res)
  if err != nil {
    log.Printf("Error decoding json: %s\n", err)
    return nil, err
  }

  return &res.Hits.Hits, nil
}

func ParseErrorResponse(response *opensearchapi.Response) (*models.SearchResponseError, error) {

  var errRes models.SearchResponseError
  err := json.NewDecoder(response.Body).Decode(&errRes)
  if err != nil {
    log.Printf("Error decoding json: %s\n", err)
    return nil, err
  }
  
  return &errRes, nil
}
