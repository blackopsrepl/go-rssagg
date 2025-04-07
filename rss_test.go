package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchFeed(t *testing.T) {
	// Test server with sample RSS response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		sampleRSS := `<?xml version="1.0" encoding="UTF-8"?>
            <rss version="2.0">
                <channel>
                    <title>AWS</title>
                    <link>https://aws.amazon.com/blogs/aws/feed/</link>
                    <description>Test Description</description>
                </channel>
            </rss>`
		w.Write([]byte(sampleRSS))
	}))
	defer testServer.Close()

	testCases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Valid RSS feed",
			url:     testServer.URL,
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			url:     "http://nonexistent.domain.xyz",
			wantErr: true,
		},
		{
			name:    "Empty URL",
			url:     "",
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchFeed(tt.url)

			// Check error conditions
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Only verify content for successful case
			if !tt.wantErr { // ← Changed condition
				if got == nil { // ← Added nil check
					t.Errorf("fetchFeed() returned nil feed for valid case")
					return
				}
				if got.Channel.Title != "AWS" {
					t.Errorf("fetchFeed() got title = %v, want %v", got.Channel.Title, "AWS")
				}
				if got.Channel.Link != "https://aws.amazon.com/blogs/aws/feed/" {
					t.Errorf("fetchFeed() got link = %v, want %v", got.Channel.Link, "https://aws.amazon.com/blogs/aws/feed/")
				}
			}
		})
	}
}

func TestFetchFeedTimeout(t *testing.T) {
	// Test server that delays response beyond timeout
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second)
		w.Write([]byte("Delayed response"))
	}))
	defer testServer.Close()

	_, err := fetchFeed(testServer.URL)
	if err == nil {
		t.Error("fetchFeed() expected timeout error, got nil")
	}
}
