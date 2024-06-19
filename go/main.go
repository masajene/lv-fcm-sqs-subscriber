package main

import (
	"context"
	"encoding/json"
	"fcm-sub/fcm"
	"fcm-sub/logger"
	"fcm-sub/model"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"strconv"
	"time"
)

var version string
var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func init() {
	version = os.Getenv("VERSION")
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var err error
	splitThreshold, err := strconv.Atoi(os.Getenv("SPLIT_THRESHOLD"))
	if err != nil {
		logger.GetLogger().Error("Error converting SPLIT_THRESHOLD to int", "error", err)
		return err
	}
	for _, message := range sqsEvent.Records {
		body := message.Body
		logger.GetLogger().Info("Processing request", "body", body)

		var payload model.SQSPayload
		err := json.Unmarshal([]byte(body), &payload)
		if err != nil {
			logger.GetLogger().Error("Error unmarshalling JSON", "error", err)
			continue
		}

		// if RegistrationIds is count over splitThreshold, return error
		if len(payload.RegistrationIDs) > splitThreshold {
			err = fmt.Errorf("RegistrationIds count over %d", splitThreshold)
			continue
		}

		// Send FCM
		err = fcm.NewFcm().Send(payload)
	}
	return err
}

func main() {
	logger.GetLogger().Info("Starting server", "jst", jst)
	logger.GetLogger().Info("Version", "version:", version)
	lambda.Start(Handler)
}
