package api

import (
	"strings"
	"testing"
)

func TestBuildShortUrl(t *testing.T) {
	type Input struct {
		host  string
		id    string
		isTLS bool
	}

	type Output struct {
		url string
	}

	tests := map[string]struct {
		input  Input
		output Output
	}{
		"Test 01 - Valid HTTPS URL": {
			Input{host: "ehgm.com.br", id: "1q2w3e4", isTLS: true},
			Output{url: "https://ehgm.com.br/r/1q2w3e4"}},

		"Test 02 - Valid HTTP URL": {
			Input{host: "ehgm.com.br", id: "1q2w3e4", isTLS: false},
			Output{url: "http://ehgm.com.br/r/1q2w3e4"}},
	}

	for i, test := range tests {
		output := buildShortUrl(test.input.host, test.input.id, test.input.isTLS)
		if output != test.output.url {
			t.Errorf("#%s: Output %v: should be: %s", i, output, test.output)
		}
	}
}

func TestValidateUrl(t *testing.T) {
	type Input struct {
		rawUrl string
	}

	type Output struct {
		hasError bool
	}

	tests := map[string]struct {
		input  Input
		output Output
	}{
		"Test 01 - Invalid URL format": {
			Input{rawUrl: "ehgm.com.br"},
			Output{hasError: true}},

		"Test 02 - Invalid URL lenght": {
			Input{rawUrl: "ehgm.com.br" + strings.Repeat("a/", 2048)},
			Output{hasError: true}},

		"Test 03 - Valid URL": {
			Input{rawUrl: "https://ehgm.com.br"},
			Output{hasError: false}},
	}

	for i, test := range tests {
		err := validateUrl(test.input.rawUrl)
		if test.output.hasError && err == nil {
			t.Errorf("#%s: Output is: %s. But should has error: %v", i, err, test.output.hasError)
			continue
		}
		if !test.output.hasError && err != nil {
			t.Errorf("#%s: Output is: %s. But should not has error: %v", i, err, test.output.hasError)
		}
	}
}
