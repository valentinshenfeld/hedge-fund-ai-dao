package main

import (
	"context"
	"testing"
)

func TestExecuteTradeHandler(t *testing.T) {
	agent := &TraderAgent{
		ctx: context.Background(),
	}

	payload := []byte("buy 1 ETH")
	resp, err := agent.ExecuteTradeHandler(payload)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(resp) != "TRADE_EXECUTED" {
		t.Errorf("Expected TRADE_EXECUTED, got %s", string(resp))
	}
}
