package model

type CompleteSqsPayload struct {
	ConnectionSource string `json:"connection_source"`
	InfoId           string `json:"info_id"`
}
