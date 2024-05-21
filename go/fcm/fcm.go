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

	br, err := client.SendEachForMulticast(context.Background(), message)
	if err != nil {
		return err
	}
	if br.FailureCount > 0 {
		var failedTokens []model.ErrorSqsPayload
		for idx, resp := range br.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedTokens = append(failedTokens, model.ErrorSqsPayload{
					ConnectionSource: m.ConnectionSource,
					InfoId:           m.Data.Message.ID,
					RegistrationId:   m.RegistrationIDs[idx],
				})
			}
		}
		// send error queue
		if len(failedTokens) > 0 {
			logger.GetLogger().Error("Failed to send message", "failedTokens", failedTokens)
			_ = sqs.SendErrorMessages(failedTokens)
		}
	}
	logger.GetLogger().Info("Successfully sent message", "response", br)
	if *m.Data.Message.DeleteFlag != "1" {
		logger.GetLogger().Info("Delete message", "message", m.Data.Message.ID)
		_ = sqs.SendCompleteMessages(m.Data.Message.ID, m.ConnectionSource)
	}
	return err
}
