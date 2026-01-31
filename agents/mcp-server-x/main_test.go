package main

import (
	"context"
	"testing"
)

func TestSearchTweetsHandler(t *testing.T) {
	resp, err := SearchTweetsHandler(context.Background(), SearchTweetsArgs{Query: "$ETH", MaxResults: 5})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data := resp.(map[string]interface{})
	if data["status"] != "success" {
		t.Errorf("Expected status success, got %s", data["status"])
	}
	if data["count"] == 0 {
		t.Errorf("Expected count > 0")
	}
}

func TestGetUserSentimentHandler(t *testing.T) {
	resp, err := GetUserSentimentHandler(context.Background(), GetUserSentimentArgs{Username: "vitalik"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data := resp.(map[string]interface{})
	if data["username"] != "vitalik" {
		t.Errorf("Expected username vitalik, got %s", data["username"])
	}
}
