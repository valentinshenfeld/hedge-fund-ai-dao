package main

import (
	"context"
	"testing"
	"github.com/google/generative-ai-go/genai"
)

type MockModel struct {
	Resp *genai.GenerateContentResponse
	Err  error
}

func (m *MockModel) GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
	return m.Resp, m.Err
}

func TestAnalyzeTokenHandler(t *testing.T) {
	mockResp := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text("Bullish sentiment for ETH"),
					},
				},
			},
		},
	}
	agent := &AnalystAgent{
		model: &MockModel{Resp: mockResp},
		ctx:   context.Background(),
	}

	payload := []byte("ETH")
	resp, err := agent.AnalyzeTokenHandler(payload)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(resp) != "Bullish sentiment for ETH" {
		t.Errorf("Expected 'Bullish sentiment for ETH', got %s", string(resp))
	}
}
