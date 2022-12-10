package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/RCGuo/aws-microservices-go/modules/checkout"
)

func wokerHandler(ctx context.Context, event events.SQSEvent) error {
  log.Println("checkoutOrderHandler")

  if event.Records != nil {
    return sqsInvocation(event)
  } else {
    log.Println("invalid event request of creating an ordering")
    return errors.New("invalid event request of creating an ordering")
  }
}

func sqsInvocation(event events.SQSEvent) error {
  log.Println("sqsInvocation")

  for _, record := range event.Records {
    log.Println("SQS record.Body: ", record.Body)

    eventBridgeEvent := events.CloudWatchEvent{}
    err := json.Unmarshal([]byte(record.Body), &eventBridgeEvent)
    if err != nil {
      log.Printf("eventBridgeEvent JSON unmarshal failed: %s\n", err)
      return err
    }
    log.Printf("eventBridgeEvent: %s", eventBridgeEvent)

    checkoutOrder := checkout.CheckoutOrder{}
    err = json.Unmarshal([]byte(eventBridgeEvent.Detail), &checkoutOrder)
    if err != nil {
      log.Printf("checkout data JSON unmarshal failed: %s\n", err)
      return err
    }

    err = checkoutOrder.DeleteUserAllCartItem()
    if err != nil {
      log.Printf("Create order error: %s\n", err.Error())
      return err
    }
  }

  return nil
}

func main(){
  lambda.Start(wokerHandler)
}