package api_utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

var headers = map[string]string{
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Headers": "*",
  "Access-Control-Allow-Methods": "*",
  "Content-Type": "application/json",
}

func APIResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	bytes, _ := json.Marshal(&body)

	return events.APIGatewayProxyResponse{
    Headers:    headers,
		Body:       string(bytes),
		StatusCode: statusCode,
	}, nil
}

func APIErrResponse(statusCode int, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
    Headers:    headers,
		Body:       err.Error(),
		StatusCode: statusCode,
	}, err
}

func APISuccessResponse(body interface{}) (events.APIGatewayProxyResponse, error) {
	bytes, _ := json.Marshal(&body)

	return events.APIGatewayProxyResponse{
    Headers:    headers,
		Body:       string(bytes),
		StatusCode: http.StatusOK,
	}, nil
}

func APIBadRequest(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
    Headers:    headers,
		Body:       err.Error(),
		StatusCode: http.StatusBadRequest,
	}, err
}

func APIServerError(err error) (events.APIGatewayProxyResponse, error) {
	printStackTrace(err)

	return events.APIGatewayProxyResponse{
    Headers:    headers,
		Body:       "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}, err
}

type stacktracer interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
}

func printStackTrace(err error) {
	var errStack errors.StackTrace

	for err != nil {
		// Find the earliest error.StackTrace
		if t, ok := err.(stacktracer); ok {
			errStack = t.StackTrace()
		}
		if c, ok := err.(causer); ok {
			err = c.Cause()
		} else {
			break
		}
	}
	if errStack != nil {
		fmt.Println(err)
		fmt.Printf("%+v\n", errStack)
	} else {
		fmt.Printf("%+v\n", errors.WithStack(err))
	}
}