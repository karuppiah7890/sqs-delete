package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/karuppiah7890/sqs-delete/pkg/config"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
)

// TODO: Write tests for all of this

var version string = "dev"

func checkSignal(signals chan os.Signal, done chan bool) {
	<-signals
	done <- true
}

func main() {
	done := make(chan bool, 1)
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt)

	go checkSignal(signals, done)

	log.Printf("version: %v", version)
	c, err := config.NewConfigFromEnvVars()
	if err != nil {
		log.Fatalf("error occurred while getting configuration from environment variables: %v", err)
	}

	awsconfig, err := awsconf.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("error occurred while loading aws configuration: %v", err)
	}

	sqsClient := sqs.NewFromConfig(awsconfig)

	queueUrl := c.GetSqsQueueUrl()

	filePath := "messages.json"

	jsonFileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error occurred while reading from file at %s: %v", filePath, err)
	}

	var messages []Message

	err = json.Unmarshal(jsonFileData, &messages)
	if err != nil {
		log.Fatalf("error occurred while deserializing messages from JSON file at %s: %v", filePath, err)
	}

	for _, message := range messages {
		err := deleteMessageFromQueue(queueUrl, sqsClient, message)
		if err != nil {
			log.Fatalf("error occurred while getting messages from sqs queue: %v", err)
		}
		fmt.Printf(".")
	}
}

type Message struct {
	// An identifier associated with the act of receiving the message. A new receipt
	// handle is returned every time you receive a message. When deleting a message,
	// you provide the last received receipt handle to delete the message.
	ReceiptHandle string `json:"receipt_handle"`
}

// Delete message from the queue
func deleteMessageFromQueue(queueUrl string, sqsClient *sqs.Client, message Message) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: &message.ReceiptHandle,
	}

	_, err := sqsClient.DeleteMessage(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("error occurred while deleting message from sqs using receipt handle '%s': %v", message.ReceiptHandle, err)
	}

	return nil
}
