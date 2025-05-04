package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/redis/go-redis/v9"
)

func saveToRedis(address string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()

	err := client.Set(ctx, address, crawler(address), 0).Err()
	if err != nil {
		log.Fatalf("Error saving to Redis: %s", err)
	}

}

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

	// Convert HTML to plain text
	plainText := htmlToPlainText(string(body))

	log.Print(string(plainText))

	return string(plainText)
}

// Simple HTML to plain text converter
func htmlToPlainText(html string) string {
	// Replace line breaks with spaces
	text := strings.ReplaceAll(html, "\n", " ")

	// Replace common HTML entities
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&nbsp;", " ")

	// Remove script and style tags and their content
	re := regexp.MustCompile("(?i)<script[^>]*>.*?</script>|<style[^>]*>.*?</style>")
	text = re.ReplaceAllString(text, "")

	// Replace paragraph and break tags with line breaks
	text = regexp.MustCompile("(?i)</p>|<br[^>]*>|<div[^>]*>").ReplaceAllString(text, "\n")

	// Add extra line breaks for headings
	text = regexp.MustCompile("(?i)</h[1-6]>").ReplaceAllString(text, "\n\n")

	// Remove all remaining HTML tags
	re = regexp.MustCompile("<[^>]*>")
	text = re.ReplaceAllString(text, "")

	// Collapse multiple spaces
	text = regexp.MustCompile(`\s{2,}`).ReplaceAllString(text, " ")

	// Collapse multiple line breaks
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	// Trim leading/trailing whitespace
	text = strings.TrimSpace(text)

	return text
}
