package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func httpException(status int, msg string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       msg,
	}, nil
}

func errCheck(msg string, err error) {
	if err != nil {
		log.Fatalf(msg)
		panic(err)
	}
}
