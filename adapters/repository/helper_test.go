package repository

import (
	"testing"
	"time"

	"ehgm.com.br/url-shortener/domain/model"
)

func TestStructToJson(t *testing.T) {
	type Input struct {
		shortUrl model.ShortUrl
	}

	type Output struct {
		json     string
		hasError bool
	}

	tests := map[string]struct {
		input  Input
		output Output
	}{
		"Test 01 - Filled struct": {
			Input{model.ShortUrl{
				Id:         "1q2w3e4r",
				Url:        "https://ehgm.com.br",
				CreateTime: time.Date(2021, 11, 15, 17, 17, 17, 0, time.Now().Location()),
				Enable:     true,
				Clicks:     10,
			}},
			Output{json: "{\"id\":\"1q2w3e4r\",\"url\":\"https://ehgm.com.br\",\"createTime\":\"2021-11-15T17:17:17.0000000-03:00\",\"enable\":true,\"clicks\":10}"},
		},

		"Test 02 - Empty struct": {
			Input{model.ShortUrl{}},
			Output{json: "{\"id\":\"\",\"url\":\"\",\"createTime\":\"0001-01-01T00:00:00Z\",\"enable\":false,\"clicks\":0}"},
		},
	}

	for i, test := range tests {
		json, err := structToJson(test.input)
		if test.output.hasError && err == nil {
			t.Errorf("#%s: Output is: %s. But should has error: %v", i, err, test.output.hasError)
			continue
		}
		if !test.output.hasError && err != nil {
			t.Errorf("#%s: Output is: %s. But should not has error: %v", i, err, test.output.hasError)
			continue
		}
		if err != nil && json == test.output.json {
			t.Errorf("#%s: Output %v: should be valid: %s", i, json, test.output.json)
		}
	}
}

func TestJsonToStruct(t *testing.T) {
	type Input struct {
		json string
	}

	type Output struct {
		shortUrl model.ShortUrl
		hasError bool
	}

	tests := map[string]struct {
		input  Input
		output Output
	}{
		"Test 01 - Filled struct": {
			Input{json: "{\"id\":\"1q2w3e4r\",\"url\":\"https://ehgm.com.br\",\"createTime\":\"2021-11-15T17:17:17.0000000-03:00\",\"enable\":true,\"clicks\":10}"},
			Output{
				shortUrl: model.ShortUrl{
					Id:         "1q2w3e4r",
					Url:        "https://ehgm.com.br",
					CreateTime: time.Date(2021, 11, 15, 17, 17, 17, 0, time.Now().Location()),
					Enable:     true,
					Clicks:     10,
				},
				hasError: false},
		},

		"Test 02 - Empty struct": {
			Input{json: "{\"id\":\"\",\"url\":\"\",\"createTime\":\"0001-01-01T00:00:00Z\",\"enable\":false,\"clicks\":0}"},
			Output{shortUrl: model.ShortUrl{}, hasError: false},
		},
	}

	for i, test := range tests {
		shortUrl, err := jsonToStruct(test.input.json)
		if test.output.hasError && err == nil {
			t.Errorf("#%s: Output is: %s. But should has error: %v", i, err, test.output.hasError)
			continue
		}
		if !test.output.hasError && err != nil {
			t.Errorf("#%s: Output is: %s. But should not has error: %v", i, err, test.output.hasError)
			continue
		}
		if err != nil && shortUrl == test.output.shortUrl {
			t.Errorf("#%s: Output %v: should be valid: %v", i, shortUrl, test.output.shortUrl)
		}
	}
}
