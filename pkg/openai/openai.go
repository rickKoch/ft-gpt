package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type OpenAI struct{}

type path string

type completionError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    *int   `json:"code"`
}

type choice struct {
	Text string `json:"text"`
}

const (
	URL     = "https://api.openai.com"
	VERSION = "v1"
)

var API_TOKEN = os.Getenv("OPENAI_API_KEY")

type CompletionRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
}

type CompletionResponse struct {
	Choices []choice         `json:"choices"`
	Error   *completionError `json:"error"`
}

const (
	COMPLETION path = "completions"
)

func (o *OpenAI) CreateCompletion(r *CompletionRequest) (string, error) {
	requestURL := fmt.Sprintf("%s/%s/%s", URL, VERSION, COMPLETION)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(r)
	if err != nil {
		fmt.Printf("error encoding json: %s\n", err)
		panic(err)
	}

	req, err := http.NewRequest("POST", requestURL, &buf)
	if err != nil {
		fmt.Printf("error creating http request: %s\n", err)
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", API_TOKEN))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error sending http request: %s\n", err)
		panic(err)
	}
	defer resp.Body.Close()

	var data CompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Printf("error: %#v\n", err)
		panic(err)
	}

	if data.Error != nil {
		panic(data.Error.Message)
	}

	return data.Choices[0].Text, nil
}
