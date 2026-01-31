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

func TestValidateRiskHandler_Pass(t *testing.T) {
	mockResp := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text(`{"status": "pass", "reason": "Safe"}`),
					},
				},
			},
		},
	}
	agent := &RiskAgent{
		model: &MockModel{Resp: mockResp},
		ctx:   context.Background(),
	}

	payload := []byte("test strategy")
	resp, err := agent.ValidateRiskHandler(payload)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(resp) != "PASS" {
		t.Errorf("Expected PASS, got %s", string(resp))
	}
}

func TestValidateRiskHandler_Fail(t *testing.T) {
	mockResp := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text(`{"status": "fail", "reason": "High slippage"}`),
					},
				},
			},
		},
	}
	agent := &RiskAgent{
		model: &MockModel{Resp: mockResp},
		ctx:   context.Background(),
	}

	payload := []byte("risky strategy")
	resp, err := agent.ValidateRiskHandler(payload)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "High slippage" {
		t.Errorf("Expected 'High slippage', got %v", err)
	}
	if resp != nil {
		t.Errorf("Expected nil response, got %v", resp)
	}
}
