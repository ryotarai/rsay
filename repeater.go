package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os/exec"
	"strings"
)

type Repeater struct {
	queueName string
	voice     string
}

func (r *Repeater) StartLoop() {
	svc := sqs.New(nil)

	log.Println("Getting Queue URL...")
	queueUrl := getQueueUrlFromName(svc, &r.queueName)
	log.Println("QueueUrl: " + *queueUrl)

	log.Println("Checking queue...")
	for true {
		messages := receiveMessages(svc, queueUrl, 20)

		for _, message := range messages {
			req := sayingRequest{}
			err := json.Unmarshal([]byte(*message.Body), &req)
			if err != nil {
				log.Print(err)
			}

			log.Println(req)

			args := []string{}
			if r.voice != "" {
				args = append(args, "-v", r.voice)
			}
			cmd := exec.Command("say", args...)
			cmd.Stdin = strings.NewReader(req.Message)
			err = cmd.Run()
			if err != nil {
				log.Print(err)
			}

			deleteMessage(svc, queueUrl, message.ReceiptHandle)
		}
	}
}

type sayingRequest struct {
	Message string
}

func getQueueUrlFromName(svc *sqs.SQS, queueName *string) *string {
	params := &sqs.GetQueueUrlInput{
		QueueName: aws.String(*queueName),
	}
	resp, err := svc.GetQueueUrl(params)

	if err != nil {
		log.Fatal(err)
	}

	return resp.QueueUrl
}

func receiveMessages(svc *sqs.SQS, queueUrl *string, waitTimeSeconds int64) []*sqs.Message {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:        aws.String(*queueUrl),
		WaitTimeSeconds: aws.Int64(waitTimeSeconds),
	}
	resp, err := svc.ReceiveMessage(params)

	if err != nil {
		log.Fatal(err)
	}

	return resp.Messages
}

func deleteMessage(svc *sqs.SQS, queueUrl *string, receiptHandle *string) {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(*queueUrl),
		ReceiptHandle: aws.String(*receiptHandle),
	}
	_, err := svc.DeleteMessage(params)

	if err != nil {
		log.Fatal(err)
	}

	return
}
