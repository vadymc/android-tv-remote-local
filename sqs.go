package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	SQS_QUEUE = "tv_remote"
)

var (
	sqsClient *sqs.SQS
	queueUrl  *string
)

func startSqs() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		panic(err)
	}
	sqsClient = sqs.New(sess)
	initQueueUrl()
	log.Println("Connected to SQS")

	pollAndProcess()
}

func initQueueUrl() {
	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(SQS_QUEUE),
	})

	if err != nil {
		log.Fatalf("Error %v", err)
		return
	}
	queueUrl = result.QueueUrl
}

func pollAndProcess() {
	chnMessages := make(chan *string, 10)
	go pollSqs(chnMessages)

	for message := range chnMessages {
		executeLiteralCommand(*message)
	}
}

func pollSqs(chn chan<- *string) {
	for {
		output, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            queueUrl,
			MaxNumberOfMessages: aws.Int64(10),
			WaitTimeSeconds:     aws.Int64(20),
		})

		if err != nil {
			log.Fatalf("failed to fetch sqs message %v", err)
			continue
		}

		for _, message := range output.Messages {
			chn <- message.Body
			deleteMessage(message)
		}
	}
}

func deleteMessage(message *sqs.Message) {
	_, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueUrl,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		log.Fatalf("Error %v", err)
	}
}
