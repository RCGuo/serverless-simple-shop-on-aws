package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridgetypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"

	"github.com/RCGuo/aws-microservices-go/datasources/dynamodb/dbclient"
	"github.com/RCGuo/aws-microservices-go/datasources/eventbridge/evclient"
	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/RCGuo/aws-microservices-go/utils/auth_utils"
)

func checkout(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("checkout")

  checkoutRequest := models.Order{}
  err := json.Unmarshal([]byte(event.Body), &checkoutRequest)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return api_utils.APIServerError(err)
  }
  
  claims, err := auth_utils.ParseCognitoClaims(event)
  if err != nil {
    log.Printf("Fetching authentication data failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  checkoutRequest.UserId = claims.Username
  if checkoutRequest.UserId == "" {
    log.Printf("user name should exist in checkout request. checkoutRequest: %+v", checkoutRequest)
    return api_utils.APIErrResponse(400, errors.New("userName should exist in checkout request"))
  }

  cartItems, gbErr := getCartItemsByUserId(checkoutRequest.UserId)
  if gbErr != nil {
    log.Printf("Getting busket error: %s\n", gbErr.Error())
    return api_utils.APIServerError(gbErr)
  }

  checkoutPayload, popErr := prepareOrderPayload(checkoutRequest, cartItems)
  if popErr != nil {
    log.Printf("Prepare order payload error: %s\n", popErr.Error())
    return api_utils.APIServerError(popErr)
  }

  publishedEvent, pcbeErr := publishCheckoutCartEvent(checkoutPayload)
  if pcbeErr != nil {
    log.Printf("Publish checkout event error: %s\n", pcbeErr.Error())
    return api_utils.APIServerError(pcbeErr)
  }

  log.Println("# publishedEvent: ", publishedEvent)

  _, dbErr := deleteCartItem(checkoutRequest.UserId)
  if dbErr != nil {
    log.Printf("Deletting busket error: %s\n", dbErr.Error())
    return api_utils.APIServerError(dbErr)
  }

  return api_utils.APISuccessResponse(nil)
}

// func roundFloat(val float64, precision uint) float64 {
// 	ratio := math.Pow(10, float64(precision))
// 	return math.Round(val*ratio) / ratio
// }

func prepareOrderPayload(checkoutRequest models.Order, carts models.CartItems) (models.Order, error) {
  log.Println("prepareOrderPayload")

  if reflect.DeepEqual(models.CartItems{}, carts) {
    return models.Order{}, errors.New("cart is empty")
  }

  var totalPrice float64
  for _, item := range carts {
    totalPrice = totalPrice + item.Price
    checkoutRequest.Items = append(checkoutRequest.Items, item)
  }
  // checkoutRequest.Total = roundFloat(totalPrice, 2)
  checkoutRequest.Total = fmt.Sprintf("%f", totalPrice)

  log.Printf("checkoutRequest : %+v", checkoutRequest)

  return checkoutRequest, nil
}

func publishCheckoutCartEvent(checkoutPayload models.Order) (*eventbridge.PutEventsOutput, error){

  checkoutPayloadJson, err := json.Marshal(checkoutPayload)
  if err != nil {
    log.Printf("JSON marshal error: %s\n", err.Error())
    return &eventbridge.PutEventsOutput{}, err
  }

  log.Printf("# checkoutPayloadJson: %s", checkoutPayloadJson)
  log.Print("# EVENT_SOURCE: ", os.Getenv("EVENT_SOURCE"))
  log.Print("# EVENT_DETAILTYPE: ", os.Getenv("EVENT_DETAILTYPE"))
  log.Print("# EVENT_BUS_NAME: ", os.Getenv("EVENT_BUS_NAME"))
  params := &eventbridge.PutEventsInput{
  	Entries: []eventbridgetypes.PutEventsRequestEntry{{
      Source:       aws.String(os.Getenv("EVENT_SOURCE")),
      Detail:       aws.String(string(checkoutPayloadJson)),
      DetailType:   aws.String(os.Getenv("EVENT_DETAILTYPE")),
      EventBusName: aws.String(os.Getenv("EVENT_BUS_NAME")),
    }},
  }

  data, err := evclient.Client.PutEvents(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("EventBridge put event API failed: %s\n", err.Error())
    return &eventbridge.PutEventsOutput{}, err
  }

  return data, nil
}

func getCartItemsByUserId(userId string) (models.CartItems, error) {

  cartItems := models.CartItems{}
  keyCond := expression.Key("userId").Equal(expression.Value(userId))
  expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		log.Printf("Failed to create dynamodb exprssion: %s", err.Error())
		return cartItems, err
  }

  params := &dynamodb.QueryInput{
    TableName: aws.String(os.Getenv("DYNAMODB_CART_TABLE_NAME")),
    ExpressionAttributeNames:  expr.Names(),
    ExpressionAttributeValues: expr.Values(),
    KeyConditionExpression:    expr.KeyCondition(),
  }
  
  data, err := dbclient.Client.Query(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("DynamoDB Query API call failed: %s", err.Error())
    return cartItems, err
  }

  err = attributevalue.UnmarshalListOfMaps(data.Items, &cartItems)
  if err != nil {
    log.Printf("Failed to unmarshal Dynamodb item, %s", err.Error())
    return cartItems, err
  }

  return cartItems, nil
}

func deleteCartItem(userId string) (string, error) {

  params := &dynamodb.DeleteItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_CART_TABLE_NAME")),
    Key: map[string]dynamodbtypes.AttributeValue{
      "userId": &dynamodbtypes.AttributeValueMemberS{Value: userId}, 
    },
  }

  deleteResult, err := dbclient.Client.DeleteItem(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("DynamoDB delete API call failed: %s", err)
    return "", err
  }

  log.Println("DynamoDB delete result: ", deleteResult)

  return "Deleted", nil
}

func main() {
  lambda.Start(checkout)
}