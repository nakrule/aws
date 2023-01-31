package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"What is your name?"`
	Age  int    `json:"How old are you?"`
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

func HandleLambdaEvent(event MyEvent) (MyResponse, error) {
	return MyResponse{Message: fmt.Sprintf("%s is %d years old! From function B.", event.Name, event.Age)}, nil
}

// "main" must be set as "handler" in "Runtime Settings" for the Lambda function.
// This setup the function main() as the entry point for Lambda.
func main() {
	lambda.Start(HandleLambdaEvent)
}
