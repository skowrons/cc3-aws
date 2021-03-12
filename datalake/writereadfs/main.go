package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// path to shared filesystem
const BasePath = "/mnt/datalake"

// Handle write file to filesystem
func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// i.e. /filesystem/path/to/file.json
	path := strings.Split(r.Path, "/filesystem")
	// paths = ["path", "to", "file.json"]
	paths := strings.Split(path[1], "/")
	directoryPath := BasePath + strings.Join(paths[:len(paths)-1], "/")
	filePath := BasePath + path[1]

	if r.HTTPMethod == http.MethodPost {
		log.Println("Create file on path:", filePath)
		log.Println("Used directory path:", directoryPath)

		// create directories
		dirErr := os.MkdirAll(directoryPath, os.ModePerm)
		if dirErr != nil {
			log.Println("ERROR creating directories: ", dirErr)
		}

		// create / override file
		f, err := os.Create(filePath)
		if err != nil {
			log.Println("ERROR while creating file: ", err)
		}
		f.WriteString(r.Body)

		data, _ := ioutil.ReadFile(filePath)
		log.Println(data)

		return events.APIGatewayProxyResponse{Body: "File written to filesystem.", StatusCode: 200}, nil
	}

	if r.HTTPMethod == http.MethodGet {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Println(err)
		}
		return events.APIGatewayProxyResponse{Body: string(data), StatusCode: 200}, nil
	}

	return events.APIGatewayProxyResponse{Body: "some unhandled method", StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
