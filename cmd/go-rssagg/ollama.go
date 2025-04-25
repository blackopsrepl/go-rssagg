package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func ollamaSendMessage(ollamaURL string, ollamaReq OllamaRequest) (*OLLamaResponse, error) {
	dat, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, err
	}
	httpClient := http.Client{}
	httpRequest, err := http.NewRequest(http.MethodPost, ollamaURL, bytes.NewReader(dat))
	if err != nil {
		return nil, err
	}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	rawBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	log.Println("Raw response:", string(rawBody))

	// Reset the body for decoding
	httpResponse.Body = io.NopCloser(bytes.NewReader(rawBody))

	ollamaResp := OLLamaResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(&ollamaResp)
	return &ollamaResp, err
}
