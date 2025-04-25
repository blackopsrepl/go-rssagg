package main

import (
	"io"
	"log"
	"net/http"
)

func crawler(address string) string {
	response, err := http.Get(address)
	if err != nil {
		log.Fatalf("Error fetching URL: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}

	log.Print(string(body))

	return string(body)
}
