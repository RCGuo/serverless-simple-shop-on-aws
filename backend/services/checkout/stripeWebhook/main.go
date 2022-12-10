package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/modules/checkout"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/lithammer/shortuuid/v4"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

type Meatadata struct {
  Items       string  `json:"items"`
  UserId      string  `json:"userId"`
  Email       string  `json:"email"`
  ShippingFee string  `json:"shippingFee"`
  Subtotal    string  `json:"subtotal"`
  Total       string  `json:"Total"`
}

func stripeWebhook(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("stripeWebhook")

  endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
  stripeEvent, err := webhook.ConstructEvent([]byte(event.Body), event.Headers["Stripe-Signature"], endpointSecret)
  if err != nil {
    log.Printf("Error verifying webhook signature: %v\n", err)
    return api_utils.APIBadRequest(err)
  }

  switch stripeEvent.Type {
  case "payment_intent.succeeded":
    log.Println("payment_intent.succeeded")

    metadata, err := getMetadataFromStripeEvent(stripeEvent)
    if err != nil {
      log.Printf("Get metadata from stripe event failed: %s\n", err)
      return api_utils.APIServerError(err)
    }

    items, err := getItemsFromMetadata(metadata)
    if err != nil {
      log.Printf("Get product items from metadata failed: %s\n", err)
      return api_utils.APIServerError(err)
    }

    chargesBytes, err := json.Marshal(stripeEvent.Data.Object["charges"])
    if err != nil {
      log.Printf("JSON unmarshal failed: %s\n", err)
      return api_utils.APIServerError(err)
    }

    var charges stripe.ChargeList
    err = json.Unmarshal(chargesBytes, &charges)
    if err != nil {
      log.Printf("JSON unmarshal failed: %s\n", err)
      return api_utils.APIServerError(err)
    }

    order := checkout.CheckoutOrder{
    	PaymentIntentId: fmt.Sprintf("%v",stripeEvent.Data.Object["id"]),
    	UserId:          metadata.UserId,
    	OrderDate:       time.Unix(stripeEvent.Created, 0).Format(time.RFC3339),
    	OrderId:         shortuuid.New(),
    	Items:           items,
    	Email:           metadata.Email,
    	Address:         fmt.Sprintf("%v",charges.Data[0].BillingDetails.Address),
    	PaymentMethod:   models.PaymentMethod{
                         Type:  fmt.Sprintf("%v", charges.Data[0].PaymentMethodDetails.Type),
                         Brand: fmt.Sprintf("%v", charges.Data[0].PaymentMethodDetails.Card.Brand),
                       },
    	PaymentStatus:   stripeEvent.Type,
    	ShippingFee:     metadata.ShippingFee,
    	Subtotal:        metadata.Subtotal,
    	Total:           metadata.Total,
    }
    log.Println("# order[payment_intent.succeeded]: ", order)
    err = order.PublishToCheckoutOrder()
    if err != nil {
      log.Printf("Create order failed: %s\n", err)
      return api_utils.APIServerError(err)
    }
  case "charge.succeeded":
    log.Println("charge.succeeded")

    metadata, err := getMetadataFromStripeEvent(stripeEvent)
    if err != nil {
      log.Printf("Get metadata from stripe event failed: %s\n", err)
      return api_utils.APIServerError(err)
    }

    items, err := getItemsFromMetadata(metadata)
    if err != nil {
      log.Printf("Get product items from metadata failed: %s\n", err)
      return api_utils.APIServerError(err)
    }

    order := checkout.CheckoutOrder{
    	PaymentIntentId: fmt.Sprintf("%v",stripeEvent.Data.Object["payment_intent"]),
    	UserId:          metadata.UserId,
      Items:           items,
      PaymentStatus:   stripeEvent.Type,
    }

    err = order.PublishToUpdateOrder()
    if err != nil {
      log.Printf("Update order failed: %s\n", err)
      return api_utils.APIServerError(err)
    }
    err = order.PublishToDeleteCart()
    if err != nil {
      log.Printf("Delete cart items failed: %s\n", err)
      return api_utils.APIServerError(err)
    }
  case "charge.failed":
    log.Println("charge.failed")

    metadata, err := getMetadataFromStripeEvent(stripeEvent)
    if err != nil {
      log.Printf("Get metadata from stripe event failed: %s\n", err)
      return api_utils.APIServerError(err)
    }

    order := checkout.CheckoutOrder{
    	PaymentIntentId: fmt.Sprintf("%v",stripeEvent.Data.Object["payment_intent"]),
    	UserId:          metadata.UserId,
      PaymentStatus:   stripeEvent.Type,
    }

    err = order.PublishToUpdateOrder()
    if err != nil {
      log.Printf("Update order failed: %s\n", err)
      return api_utils.APIServerError(err)
    }
  default:
    log.Printf("Unhandled event type: %s\n", stripeEvent.Type)
  }

  return api_utils.APISuccessResponse(nil)
}

func getMetadataFromStripeEvent(stripeEvent stripe.Event) (Meatadata, error) {
  var metadata Meatadata
  bytes, err := json.Marshal(stripeEvent.Data.Object["metadata"])
  if err != nil {
    log.Printf("JSON Marshal failed: %s\n", err)
    return metadata, err
  }
  
  err = json.Unmarshal([]byte(bytes), &metadata)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return metadata, err
  }
  log.Println("### metadata to  req:  ", metadata)

  return metadata, nil
}

func getItemsFromMetadata(metadata Meatadata) (models.CartItems, error) {
  var items models.CartItems
  err := json.Unmarshal([]byte(metadata.Items), &items)
  if err != nil {
    log.Printf("JSON unmarshal failed: %s\n", err)
    return items, err
  }
  return items, nil
}

func main() {
  lambda.Start(stripeWebhook)
}
