package main

import (
	"testing"
)

func TestCrawler(t *testing.T) {
	const url string = "https://aws.amazon.com/blogs/aws/aws-announces-pixtral-large-25-02-model-in-amazon-bedrock-serverless/"
	crawler(url)
}
