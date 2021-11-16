package repository

import (
	"encoding/json"

	"ehgm.com.br/url-shortener/domain/model"
)

func structToJson(object interface{}) (string, error) {
	text, err := json.Marshal(object)
	return string(text), err
}

func jsonToStruct(text string) (model.ShortUrl, error) {
	var object model.ShortUrl
	err := json.Unmarshal([]byte(text), &object)
	return object, err
}
