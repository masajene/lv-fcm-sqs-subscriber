package model

type ErrorSqsPayload struct {
	ConnectionSource string `json:"connection_source"`
	InfoId           string `json:"info_id"`
	RegistrationId   string `json:"registration_id"`
}
