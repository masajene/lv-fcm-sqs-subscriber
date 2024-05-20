package model

import (
	"encoding/json"
)

// SQSPayload struct defines the expected payload for SQS messages
type SQSPayload struct {
	ConnectionSource string         `json:"connection_source"`
	RegistrationIDs  []string       `json:"registration_ids" binding:"required"`
	TimeToLive       int            `json:"time_to_live"`
	CollapseKey      string         `json:"collapse_key" binding:"required"`
	Data             SQSPayloadData `json:"data" binding:"required"`
	Priority         string         `json:"priority" binding:"required"`
}

type SQSPayloadData struct {
	Message SQSPayloadMessage `json:"message"`
}

type SQSPayloadMessage struct {
	Type               string           `json:"type" binding:"required"`
	APIURL             string           `json:"api_url" binding:"required"`
	ID                 string           `json:"id" binding:"required"`
	Revision           *string          `json:"revision" binding:"required"`
	Title              *string          `json:"title" binding:"required"`
	DeliveryDatetime   *string          `json:"delivery_datetime" binding:"required"`
	EmergencyLevel     *string          `json:"emergency_level" binding:"required"`
	CategoryID         *string          `json:"category_id" binding:"required"`
	Registered         *string          `json:"registered"`
	DeleteFlag         *string          `json:"delete_flag" binding:"required"`
	Content            *string          `json:"content"`
	InfoIDs            *[]string        `json:"info_ids"`
	EmgIDs             *[]string        `json:"emg_ids"`
	AudioData1         *string          `json:"audio_data1"`
	AudioData2         *string          `json:"audio_data2"`
	AudioData3         *string          `json:"audio_data3"`
	Image              *SQSPayloadImage `json:"image"`
	Sound              *SQSPayloadSound `json:"sound"`
	ExternalLinkURL    *string          `json:"external_link_url"`
	ExternalLinkTitle  *string          `json:"external_link_title"`
	EmergencyMode      bool             `json:"emergency_mode"`
	IncompleteFlag     bool             `json:"incomplete_flag"`
	ExpireDatetime     string           `json:"expire_datetime"`
	ToBeConfirmed      int              `json:"to_be_confirmed"`
	EmergencyModeTitle string           `json:"emergency_mode_title"`
	Critical           int              `json:"critical"`
}

type SQSPayloadImage struct {
	ID            string `json:"id"`
	ImageScalable bool   `json:"image_scalable"`
	ImageURL      string `json:"image_url"`
	ImageThumbURL string `json:"image_thumb_url"`
}

type SQSPayloadSound struct {
	ID        string `json:"id"`
	SoundURL  string `json:"sound_url"`
	DelayTime int    `json:"delay_time"`
}

func (d SQSPayloadData) ConvertFcmPayload() (map[string]string, error) {
	messageJSON, err := json.Marshal(d.Message)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"message": string(messageJSON),
	}, nil
}
