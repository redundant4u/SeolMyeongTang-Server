package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/redundant4u/SeolMyeongTang-Server/internal"
)

func main() {
	log.Default().Println("Starting lambda")
	lambda.Start(internal.Router)
}
