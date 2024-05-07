package main

import (
	"context"
	"fcm-sub/logger"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"time"
)

var version string
var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func init() {
	version = os.Getenv("VERSION")
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		body := message.Body
		logger.GetLogger().Info("Processing request", "body", body)
	}
	return nil
}

func main() {
	logger.GetLogger().Info("Starting server", "jst", jst)
	logger.GetLogger().Info("Version", "version:", version)
	lambda.Start(Handler)
}
