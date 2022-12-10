package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/RCGuo/aws-microservices-go/models"
	"github.com/RCGuo/aws-microservices-go/utils/api_utils"
	"github.com/RCGuo/aws-microservices-go/utils/auth_utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

func calculateOrderAmount(items models.CartItems) int64 {
  var totalPrice float64
  for _, item := range items {
    totalPrice = totalPrice + item.Price * float64(item.Quantity) * 100
  }

  return int64(totalPrice)
}

func createPaymentIntent(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Println("createPaymentIntent")

  if event.HTTPMethod != "POST" {
    return api_utils.APIResponse(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
  }

  claims, err := auth_utils.ParseCognitoClaims(event)
  if err != nil {
    log.Printf("Fetching authentication data failed: %s\n", err)
    return api_utils.APIServerError(err)
  }

  var req struct {
    Items models.CartItems `json:"items"`
  }
  err = json.Unmarshal([]byte(event.Body), &req)
  if err != nil {
    return api_utils.APIServerError(err)
  }

  if len(req.Items) == 0 {
    return api_utils.APIBadRequest(errors.New("bad request, product item should not be empty"))
  }

  cartItems, err := json.Marshal(req.Items)
  if err != nil {
    return api_utils.APIServerError(err)
  }

  var shippingFee int64 = 0
  subTotal := calculateOrderAmount(req.Items)
  total := subTotal + shippingFee
  // Create a PaymentIntent with amount and currency
  params := &stripe.PaymentIntentParams{
    Amount:   stripe.Int64(total),
    Currency: stripe.String(string(stripe.CurrencyUSD)),
    AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
      Enabled: stripe.Bool(true),
    },
  }
  params.AddMetadata("userId", claims.Username)
  params.AddMetadata("email", claims.Email)
  params.AddMetadata("items", string(cartItems))
  params.AddMetadata("subtotal", strconv.FormatInt(subTotal, 10))
  params.AddMetadata("shippingFee", strconv.FormatInt(shippingFee, 10))
  params.AddMetadata("total", strconv.FormatInt(total, 10))
  
  pi, err := paymentintent.New(params)
  if err != nil {
    log.Printf("pi.New: %v", err)
    return api_utils.APIServerError(err)
  }

  return api_utils.APISuccessResponse(struct {
    ClientSecret string `json:"clientSecret"`
    Total        int64  `json:"total"`
  }{
    ClientSecret: pi.ClientSecret,
    Total: total,
  })
}

func main() {
  stripe.Key = os.Getenv("STRIPE_PRIVATE_KEY")
  lambda.Start(createPaymentIntent)
}
