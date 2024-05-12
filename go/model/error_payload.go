package model

type ErrorSqsPayload struct {
	InfoId         string `json:"info_id"`
	RegistrationId string `json:"registration_id"`
}
