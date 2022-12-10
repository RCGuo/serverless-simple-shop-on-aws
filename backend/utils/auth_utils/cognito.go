package auth_utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type AWSCognitoClaims struct {
  Subject  string `json:"sub"`
  Username string `json:"cognito:username"`
  Email    string `json:"email"`
  Groups   string `json:"cognito:groups"`
}

func ParseCognitoClaims(event events.APIGatewayProxyRequest) (AWSCognitoClaims , error) {
  var claims AWSCognitoClaims 
  bytes, err := json.Marshal(event.RequestContext.Authorizer["claims"])
  if err != nil {
    log.Printf("JSON Marshal failed: %s\n", err)
    return claims, err
  }
  
  err = json.Unmarshal([]byte(bytes), &claims)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return claims, err
  }

  return claims, nil
}

func IsRequestByCognitoAdmin(event events.APIGatewayProxyRequest) (bool, error) {
  claims, err := ParseCognitoClaims(event)
  if err != nil {
    return false, errors.New(fmt.Sprintf("Fetching authentication data failed: %s", err))
  }

  if claims.Groups != "admin" {
    return false, errors.New(fmt.Sprintf("Insufficient permissions: %s", err))
  }

  return true, nil
}