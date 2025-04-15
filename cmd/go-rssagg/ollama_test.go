package main

import (
	"log"
	"testing"
)

func TestOllamaSendMessage(t *testing.T) {
	msg := OllamaMessage{
		Role:    "user",
		Content: "Why is the sky green?",
	}
	req := OllamaRequest{
		Model:    "deepseek-v2:16b",
		Messages: []OllamaMessage{msg},
	}
	resp, err := ollamaSendMessage("http://localhost:11434/api/chat", req)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(resp.Message.Content)
}
