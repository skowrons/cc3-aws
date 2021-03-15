package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// PhoneNumber must be changed
const PhoneNumber = "+49 176 12312312"

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Received a github event!")
	log.Println(r.Body)

	if strings.Contains(r.Body, "github") {
		log.Println("creating session")
		sess := session.Must(session.NewSession())
		log.Println("session created")

		svc := sns.New(sess)
		log.Println("service created")
		msg := fmt.Sprintf("Achtung jemand hat gepusht!")

		params := &sns.PublishInput{
			Message:     aws.String(msg),
			PhoneNumber: aws.String(PhoneNumber),
		}
		resp, err := svc.Publish(params)
		if err != nil {
			log.Println(err.Error())
			return events.APIGatewayProxyResponse{Body: "error while trying to send sms", StatusCode: http.StatusInternalServerError}, nil
		}
		log.Println(resp)

		return events.APIGatewayProxyResponse{Body: "sms will be send", StatusCode: http.StatusOK}, nil
	}

	return events.APIGatewayProxyResponse{Body: "sms will not be send", StatusCode: http.StatusBadRequest}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
