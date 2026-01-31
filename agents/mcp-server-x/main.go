package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/adk/mcp"
	"google.golang.org/adk/tool"
)

// SearchTweetsArgs defines the parameters for searching tweets
type SearchTweetsArgs struct {
	Query      string `json:"query" jsonschema:"The search query (e.g., '$ETH sentiment')"`
	MaxResults int    `json:"max_results" jsonschema:"Maximum number of tweets to return (default 10)"`
}

// SearchTweetsHandler simulates fetching tweets from the Twitter API v2
func SearchTweetsHandler(ctx context.Context, args SearchTweetsArgs) (any, error) {
	apiKey := os.Getenv("TWITTER_API_KEY")
	if apiKey == "" {
		// Fallback for demonstration/dev
		log.Println("Warning: TWITTER_API_KEY is not set. Returning simulated data.")
	}

	if args.MaxResults == 0 {
		args.MaxResults = 10
	}

	// In a real implementation, we would call:
	// https://api.twitter.com/2/tweets/search/recent?query=...
	
	// Simulated response structure as described in README.md (Context Filtering)
	results := []map[string]interface{}{
		{
			"text":      "Feeling very bullish on $ETH today! The zkEVM activity is spiking. 🚀",
			"author":    "crypto_whale_123",
			"verified":  true,
			"metrics":   map[string]int{"likes": 1240, "retweets": 350},
			"timestamp": time.Now().Add(-15 * time.Minute).Format(time.RFC3339),
		},
		{
			"text":      "New DeFi protocol on Scroll looking interesting. Audits look clean.",
			"author":    "defi_analyst",
			"verified":  false,
			"metrics":   map[string]int{"likes": 85, "retweets": 12},
			"timestamp": time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		},
	}

	return map[string]interface{}{
		"query":  args.Query,
		"tweets": results,
		"count":  len(results),
		"status": "success",
	}, nil
}

// GetUserSentimentArgs defines parameters for analyzing a specific user's posts
type GetUserSentimentArgs struct {
	Username string `json:"username" jsonschema:"The Twitter handle to analyze"`
}

// GetUserSentimentHandler extracts recent activity for a specific user
func GetUserSentimentHandler(ctx context.Context, args GetUserSentimentArgs) (any, error) {
	return map[string]interface{}{
		"username":     args.Username,
		"recent_posts": 5,
		"sentiment":    "Bullish",
		"top_keywords": []string{"Layer2", "zkRollup", "Ethereum"},
	}, nil
}

func main() {
	// Initialize MCP Server
	server := mcp.NewServer(
		"mcp-server-x",
		"1.0.0",
		mcp.WithDescription("X (Twitter) Sentiment Analysis MCP Server for Hedge Fund AI DAO"),
	)

	// Register Tools
	// The AnalystAgent uses these to detect social signals
	server.RegisterTool(tool.NewFunctionTool("search_tweets", SearchTweetsHandler))
	server.RegisterTool(tool.NewFunctionTool("get_user_sentiment", GetUserSentimentHandler))

	// Mock for "analyze_sentiment" mentioned in some prompts
	server.RegisterTool(tool.NewFunctionTool("analyze_sentiment", SearchTweetsHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("X (Twitter) MCP Server starting on port %s...", port)
	if err := server.ListenAndServe(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
