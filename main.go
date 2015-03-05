package main

import (
	"encoding/json"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/sqs"
	"log"
	"os"
	"os/exec"
)

type SayingRequest struct {
	Message string
}

func main() {
	creds, err := aws.EnvCreds()
	if err != nil {
		log.Fatal(err)
	}

	client := sqs.New(creds, "ap-northeast-1", nil)

	if len(os.Args) < 2 {
		log.Fatal("Usage: rsay QUEUE_NAME")
	}

	queueName := os.Args[1]
	getQueueURLReq := sqs.GetQueueURLRequest{
		QueueName: &queueName,
	}

	getQueueRes, err := client.GetQueueURL(&getQueueURLReq)
	if err != nil {
		log.Fatal(err)
	}

	queueURL := getQueueRes.QueueURL
	log.Println(*queueURL)

	for true {
		log.Println("Checking queue...")
		waitTimeSeconds := 20
		receiveMessageReq := sqs.ReceiveMessageRequest{
			QueueURL:        queueURL,
			WaitTimeSeconds: &waitTimeSeconds,
		}
		receiveMessageRes, err := client.ReceiveMessage(&receiveMessageReq)
		if err != nil {
			log.Fatal(err)
		}

		for _, message := range receiveMessageRes.Messages {
			sayingRequest := SayingRequest{}
			err = json.Unmarshal([]byte(*message.Body), &sayingRequest)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(sayingRequest)

			cmd := exec.Command("say", sayingRequest.Message)
			err = cmd.Run()
			if err != nil {
				log.Fatal(err)
			}

			deleteMessageReq := sqs.DeleteMessageRequest{
				QueueURL:      queueURL,
				ReceiptHandle: message.ReceiptHandle,
			}
			client.DeleteMessage(&deleteMessageReq)
		}
	}
}
