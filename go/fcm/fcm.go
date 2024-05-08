package fcm

import (
	"context"
	_ "embed"
	"errors"
	"fcm-sub/logger"
	"fcm-sub/model"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const firebaseScope = "https://www.googleapis.com/auth/firebase.messaging"
const fcmEndpoint = "https://fcm.googleapis.com/v1/{parent=projects/*}/messages:send"

//go:embed serviceAccountKey.json
var key []byte

func NewFcm() *Fcm {
	return &Fcm{}
}

type Fcm struct {
}

func (f *Fcm) Send(m model.SQSPayload) error {
	opt := option.WithCredentialsJSON(key)
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
		Data:   dataMap,
		Android: &messaging.AndroidConfig{
			CollapseKey: m.CollapseKey,
			Priority:    m.Priority,
			Data:        dataMap,
		},
	}

	br, err := client.SendEachForMulticastDryRun(context.Background(), message)
	if err != nil {
		return err
	}
	if br.FailureCount > 0 {
		var failedTokens map[string]error
		failedTokens = make(map[string]error)
		for idx, resp := range br.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedTokens[m.RegistrationIDs[idx]] = resp.Error
			}
		}
		fmt.Printf("Map of tokens that caused failures: %v\n", failedTokens)
		// send error queue
	}
	return err
}

// NewToken function to get token for fcm-send
func (f *Fcm) newToken() (string, error) {
	cfg, err := google.JWTConfigFromJSON(key, firebaseScope)
	if err != nil {
		return "", errors.New("fcm: failed to get JWT config for the firebase.messaging scope")
	}
	ts := cfg.TokenSource(context.Background())
	token, err := ts.Token()
	if err != nil {
		return "", errors.New("fcm: failed to generate Bearer token")
	}
	return token.AccessToken, nil
}
