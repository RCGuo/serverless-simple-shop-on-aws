package checkout

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

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
)

type CheckoutOrder models.Order

func (c CheckoutOrder) PublishToCheckoutOrder() error {
  log.Println("PublishToCheckoutOrder")

  if c.UserId == "" {
    log.Printf("username must be set: %+v", c)
    return errors.New("username must be set in checkout request")
  }
  _, err := c.publishEvent(
    os.Getenv("EVENT_SOURCE"), 
    os.Getenv("CHECKOUT_EVENT_DETAILTYPE"), 
    os.Getenv("EVENT_BUS_NAME"))
  if err != nil {
    log.Printf("Publish '%s' event error: %s\n", os.Getenv("CHECKOUT_EVENT_DETAILTYPE"), err.Error())
    return err
  }
  return nil
}

func (c CheckoutOrder) PublishToUpdateOrder() error {
  log.Println("ToUpdateOrder task")

  if c.UserId == "" {
    log.Printf("username must be set: %+v", c)
    return errors.New("username must be set in checkout request")
  }
  _, err := c.publishEvent(
    os.Getenv("EVENT_SOURCE"), 
    os.Getenv("UPDATE_STATUS_EVENT_DETAILTYPE"), 
    os.Getenv("EVENT_BUS_NAME"))
  if err != nil {
    log.Printf("Publish '%s' event error: %s\n", os.Getenv("UPDATE_STATUS_EVENT_DETAILTYPE"), err.Error())
    return err
  }
  return nil
}

func (c CheckoutOrder) PublishToDeleteCart() error {
  log.Println("PublishToDeleteCart")

  if c.UserId == "" {
    log.Printf("username must be set: %+v", c)
    return errors.New("username must be set in checkout request")
  }
  _, err := c.publishEvent(
    os.Getenv("EVENT_SOURCE"), 
    os.Getenv("DELETE_CART_EVENT_DETAILTYPE"), 
    os.Getenv("EVENT_BUS_NAME"))
  if err != nil {
    log.Printf("Publish '%s' event error: %s\n", os.Getenv("DELETE_CART_EVENT_DETAILTYPE"), err.Error())
    return err
  }
  return nil
}

func (c CheckoutOrder) publishEvent(eventSource, eventDetailType, eventBusName string) (*eventbridge.PutEventsOutput, error){
  
  payload, err := json.Marshal(c)
  if err != nil {
    log.Printf("JSON marshal error: %s\n", err.Error())
    return &eventbridge.PutEventsOutput{}, err
  }

  params := &eventbridge.PutEventsInput{
  	Entries: []eventbridgetypes.PutEventsRequestEntry{{
      Source:       aws.String(eventSource),
      Detail:       aws.String(string(payload)),
      DetailType:   aws.String(eventDetailType),
      EventBusName: aws.String(eventBusName),
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

func (c CheckoutOrder) CreateOrder() (*dynamodb.PutItemOutput, error) {
  log.Println("CreateOrder")

  data, err := attributevalue.MarshalMap(c)
	if err != nil {
		log.Printf("DynamoDB attributevalue MarshalMap error: %s\n", err.Error())
    return &dynamodb.PutItemOutput{}, err
	}

  params := &dynamodb.PutItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_ORDER_TABLE_NAME")),
    Item: data,
  }

  createResult, err := dbclient.Client.PutItem(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("Create API call failed: %s", err)
    return &dynamodb.PutItemOutput{}, err
  }

  return createResult, nil
}

func (c CheckoutOrder) UpdateOrderStatus() error {
  log.Println("UpdateOrderStatus")

  update := expression.UpdateBuilder{}.Set(
    expression.Name("paymentStatus"),
    expression.Value(c.PaymentStatus))
  condition := expression.AttributeExists(expression.Name("paymentStatus"))
  expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(condition).Build()
  if err != nil {
		log.Printf("Failed to create dynamodb exprssion: %s", err.Error())
		return err
  }

  params := &dynamodb.UpdateItemInput{
    TableName: aws.String(os.Getenv("DYNAMODB_ORDER_TABLE_NAME")),
    Key: map[string]dynamodbtypes.AttributeValue{
      "userId": &dynamodbtypes.AttributeValueMemberS{Value: c.UserId}, 
      "paymentIntentId": &dynamodbtypes.AttributeValueMemberS{Value: c.PaymentIntentId}, 
    },
    ExpressionAttributeNames:  expr.Names(),
    ExpressionAttributeValues: expr.Values(),
    UpdateExpression:          expr.Update(),
    ConditionExpression:       expr.Condition(),
    ReturnValues: "UPDATED_NEW",
  }

  _, err = dbclient.Client.UpdateItem(
    context.TODO(),
    params,
  )
  if err != nil {
    log.Printf("Create API call failed: %s", err)
    return err
  }

  return nil
}

func (c CheckoutOrder) DeleteUserAllCartItem() error {
  log.Println("DeleteUserAllCartItem")

  var itemLength int = len(c.Items)
  var writeReqs  []dynamodbtypes.WriteRequest 
  for i, item := range c.Items {
    writeReqs = append(writeReqs, dynamodbtypes.WriteRequest {
      DeleteRequest: &dynamodbtypes.DeleteRequest{Key: map[string]dynamodbtypes.AttributeValue{
        "userId": &dynamodbtypes.AttributeValueMemberS{Value: c.UserId}, 
        "productId": &dynamodbtypes.AttributeValueMemberS{Value: item.ProductId}, 
      },},
    })

    if(len(writeReqs) == 25 || i+1 == itemLength) {
      params := &dynamodb.BatchWriteItemInput{
        RequestItems: map[string][]dynamodbtypes.WriteRequest{
          os.Getenv("DYNAMODB_CART_TABLE_NAME"): writeReqs,
        },
      }
    
      _, err := dbclient.Client.BatchWriteItem(
        context.TODO(),
        params,
      )
  
      if err != nil {
        log.Printf("Create API call failed: %s", err)
        return err
      }

      writeReqs = []dynamodbtypes.WriteRequest{}
    }
  }

  return nil
}