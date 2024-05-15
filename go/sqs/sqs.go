package sqs

import (
	"encoding/json"
	"fcm-sub/logger"
	"fcm-sub/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"os"
	"strings"
)

func SendErrorMessages(m []model.ErrorSqsPayload) error {
	queueURL := os.Getenv("QUEUE_URL")
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
	sqsSvc := sqs.New(sess)

	for _, v := range m {
		msgJson, err := json.Marshal(v)
		if err != nil {
			return err
		}

		uuidWithHyphen := uuid.New()
		gid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

		_, err = sqsSvc.SendMessage(&sqs.SendMessageInput{
			MessageBody:            aws.String(string(msgJson)),
			QueueUrl:               &queueURL,
			MessageGroupId:         aws.String(gid),
			MessageDeduplicationId: aws.String("dp-" + gid),
		})
		if err != nil {
			return err
		}
		logger.GetLogger().Info("Send message to SQS", "message", v)
	}
	return nil
}

func SendCompleteMessages(id, source string) error {
	queueURL := os.Getenv("COMPLETE_QUEUE_URL")
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
	sqsSvc := sqs.New(sess)

	m := model.CompleteSqsPayload{
		ConnectionSource: source,
		InfoId:           id,
	}
	msgJson, err := json.Marshal(m)
	if err != nil {
		return err
	}

	uuidWithHyphen := uuid.New()
	gid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	payload := &sqs.SendMessageInput{
		MessageBody:            aws.String(string(msgJson)),
		QueueUrl:               &queueURL,
		MessageGroupId:         aws.String(gid),
		MessageDeduplicationId: aws.String("dp-" + gid),
	}

	_, err = sqsSvc.SendMessage(payload)
	if err != nil {
		logger.GetLogger().Error("Error sending message to SQS", "error", err)
	}
	logger.GetLogger().Info("Send message to SQS", "message", payload)
	return err
}
