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
	"time"
)

var version string
var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func init() {
	version = os.Getenv("VERSION")
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var err error
	for _, message := range sqsEvent.Records {
		body := message.Body
		logger.GetLogger().Info("Processing request", "body", body)

		b2 := "{\"registration_ids\": [\"ID1\",\"ID2\",\"ID3\"],\"time_to_live\": 3600,\"collapse_key\": \"info_message_6503\",\"data\": {\"message\": {\"info_ids\": [\"6503\"],\"emg_ids\": null,\"type\": \"NORMAL\",\"api_url\": \"https://test01-api.testing.lifevision.net/api/\",\"id\": \"6503\",\"revision\": \"1\",\"title\": \"配信動作確認001 \",\"delivery_datetime\": \"2024-03-02 09:18:45\",\"emergency_level\": \"1\",\"category_id\": \"1\",\"registered\": \"テスタ\",\"delete_\": \"0\",\"audio_data1\": \"M01 新着通知音 1\",\"audio_data2\": null,\"audio_data3\": null,\"content\": \"配信動作確認です。\",\"image\": {\"id\": \"1458\",\"image_scalable\": true,\"image_url\": \"https://s3-ap-northeast-1.amazonaws.com/lifevision.testing/test01/images/17093386155192.jpg\",\"image_thumb_url\": \"https://s3-ap-northeast-1.amazonaws.com/lifevision.testing/test01/thumbs/17093386155192.jpg\"},\"sound\": {\"id\": \"1459\",\"sound_url\": \"https://s3-ap-northeast-1.amazonaws.com/lifevision.testing/test01/sounds/vot4800l0gk1c04oeoufubnjv6.mp3\",\"delay_time\": 0},\"external_link_url\": \"https://s3-ap-northeast-1.amazonaws.com/lifevision.testing/test01/pdf/1709338697167.pdf\",\"external_link_title\": \"デジタル庁HP\",\"emergency_mode\": false,\"incomplete_flag\": false}},\"priority\": \"high\"}"

		if body == b2 {
			logger.GetLogger().Info("Same body")
		} else {
			logger.GetLogger().Info("Different body")
		}

		var payload model.SQSPayload
		err := json.Unmarshal([]byte(body), &payload)
		if err != nil {
			logger.GetLogger().Error("Error unmarshalling JSON", "error", err)
			continue
		}

		// if RegistrationIds is count over 500, return error
		if len(payload.RegistrationIDs) > 500 {
			err = fmt.Errorf("RegistrationIds count over 1000")
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
