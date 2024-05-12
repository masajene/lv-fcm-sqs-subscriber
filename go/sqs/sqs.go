package sqs

import (
	"encoding/json"
	"fcm-sub/logger"
	"fcm-sub/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
	"strconv"
)

func SendErrorMessages(m []model.ErrorSqsPayload) error {
	queueURL := os.Getenv("QUEUE_URL")
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
	sqsSvc := sqs.New(sess)

	for i, v := range m {
		msgJson, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, err = sqsSvc.SendMessage(&sqs.SendMessageInput{
			MessageBody:    aws.String(string(msgJson)),
			QueueUrl:       &queueURL,
			MessageGroupId: aws.String(v.InfoId + "-" + strconv.Itoa(i)),
		})
		if err != nil {
			return err
		}
		logger.GetLogger().Info("Send message to SQS", "message", v)
	}
	return nil
}

func SendCompleteMessages(id string) error {
	queueURL := os.Getenv("COMPLETE_QUEUE_URL")
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
	sqsSvc := sqs.New(sess)

	_, err := sqsSvc.SendMessage(&sqs.SendMessageInput{
		MessageBody:    aws.String(id),
		QueueUrl:       &queueURL,
		MessageGroupId: aws.String(id),
	})
	if err != nil {
		logger.GetLogger().Error("Error sending message to SQS", "error", err)
	}
	return err
}
