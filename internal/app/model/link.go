package model

import (
	"encoding/json"
)

type Link struct {
	OriginalAddress string `json:"address"`
	ExpirationTime  int    `json:"expiration_time"`
	CreatedAt       string `json:"created_at"`
	ShortLink       string `json:"short_link"`
}

// func (sl *Link) Validate() error {

// }

func (sl *Link) ParseStringIntoStruct(str string) error {
	in := []byte(str)
	var raw map[string]interface{}
	if err := json.Unmarshal(in, &raw); err != nil {
		return err
	}

	sl.OriginalAddress = raw["address"].(string)
	sl.CreatedAt = raw["created_at"].(string)

	val := int(raw["expiration_time"].(float64))
	sl.ExpirationTime = val

	return nil
}
