package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/blackopsrepl/go-rssagg/internal/database"
	"github.com/redis/go-redis/v9"
)

// starts webcrawler worker process
// takes a db connection, a number of conturrency units and time between requests
func startCrawling(db *database.Queries, rdb *redis.Client, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Collecting page contents every %s on %v goroutines...", timeBetweenRequest, concurrency)

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		posts, err := db.GetNextPostsToCrawl(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Couldn't get next posts to crawl: %v", err)
			continue
		}
		log.Printf("Found %d posts to crawl", len(posts))

		// create a waitgroup and crawl concurrently
		wg := &sync.WaitGroup{}
		for _, post := range posts {
			wg.Add(1)
			go crawlPost(db, rdb, wg, post)
		}
		wg.Wait()
	}
}

// returns nothing, function is called concurrently from a go routine
func crawlPost(db *database.Queries, rdb *redis.Client, wg *sync.WaitGroup, post database.Post) {
	// defer and mark as crawled
	defer wg.Done()

	_, err := db.MarkPostCrawled(context.Background(), post.ID)
	if err != nil {
		log.Printf("Error marking post %s as crawled: %v", post.Url, err)
		return
	}

	// crawl post
	postData, err := crawler(post.Url)
	if err != nil {
		log.Printf("Error crawling post %s: %v", post.Url, err)
		return
	}

	// save to redis
	err = rdb.Set(context.Background(), post.Url, postData, 0).Err()
	if err != nil {
		log.Fatalf("Error saving to Redis: %s", err)
	}
}

func crawler(address string) (string, error) {
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

	return string(plainText), err
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
