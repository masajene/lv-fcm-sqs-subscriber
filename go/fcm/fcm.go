package fcm

import (
	"context"
	_ "embed"
	"fcm-sub/logger"
	"fcm-sub/model"
	"fcm-sub/sqs"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

//go:embed default-serviceAccountKey.json
var defaultKey []byte

//go:embed kariya-serviceAccountKey.json
var kariyaKey []byte

//go:embed yamato-serviceAccountKey.json
var yamatoKey []byte

func NewFcm() *Fcm {
	return &Fcm{}
}

type Fcm struct {
}

func (f *Fcm) Send(m model.SQSPayload) error {
	var opt option.ClientOption
	switch m.ConnectionSource {
	case "kariya.lifevision.net":
		opt = option.WithCredentialsJSON(kariyaKey)
	case "yamato.lifevision.net":
		opt = option.WithCredentialsJSON(yamatoKey)
	default:
		opt = option.WithCredentialsJSON(defaultKey)
	}
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}
	client, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	dataMap, err := m.Data.ConvertFcmPayload()
	if err != nil {
		return err
	}
	logger.GetLogger().Info("DataMap", "dataMap", dataMap)

	message := &messaging.MulticastMessage{
		Tokens: m.RegistrationIDs,
		//Data:   dataMap,
		Android: &messaging.AndroidConfig{
			CollapseKey: m.CollapseKey,
			Priority:    m.Priority,
			Data:        dataMap,
		},
	}

	type FcmSendResponse struct {
		Success   bool
		MessageID string
		Error     string
	}

	type FcmRequest struct {
		SuccessCount int
		FailureCount int
		Responses    []*FcmSendResponse
	}

	br, err := client.SendEachForMulticast(context.Background(), message)
	if err != nil {
		return err
	}
	msg, _ := m.Data.UnmarshalMessage()
	var rest []*FcmSendResponse

	if br.FailureCount > 0 {
		var failedTokens []model.ErrorSqsPayload
		for idx, resp := range br.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedTokens = append(failedTokens, model.ErrorSqsPayload{
					ConnectionSource: m.ConnectionSource,
					InfoId:           msg["id"],
					RegistrationId:   m.RegistrationIDs[idx],
				})
			}
			var errMsg string
			if resp.Error != nil {
				errMsg = resp.Error.Error()
			}
			var messageID string
			if idx < len(m.RegistrationIDs) {
				messageID = m.RegistrationIDs[idx]
			}
			rest = append(rest, &FcmSendResponse{
				Success:   resp.Success,
				MessageID: messageID,
				Error:     errMsg,
			})
		}
		// send error queue
		if len(failedTokens) > 0 {
			logger.GetLogger().Error("Failed to send message", "failedTokens", failedTokens)
			_ = sqs.SendErrorMessages(failedTokens)
		}
	}

	result := &FcmRequest{
		SuccessCount: br.SuccessCount,
		FailureCount: br.FailureCount,
		Responses:    rest,
	}
	logger.GetLogger().Info("Successfully sent message", "response", result)
	if msg["delete_flag"] != "1" {
		logger.GetLogger().Info("Delete message", "message", msg["id"])
		_ = sqs.SendCompleteMessages(msg["id"], m.ConnectionSource)
	}
	return err
}
