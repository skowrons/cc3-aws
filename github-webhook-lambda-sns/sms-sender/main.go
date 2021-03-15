package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)



func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("HEllow SMS")
	log.Println(r.Body)

	return events.APIGatewayProxyResponse{Body: "some unhandled method", StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
