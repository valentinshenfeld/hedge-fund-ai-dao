package main

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type MockETHClient struct {
	Balance    *big.Int
	CallResult []byte
	Header     *types.Header
	Logs       []types.Log
	Err        error
}

func (m *MockETHClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return m.Balance, m.Err
}

func (m *MockETHClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return m.CallResult, m.Err
}

func (m *MockETHClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return m.Header, m.Err
}

func (m *MockETHClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return m.Logs, m.Err
}

func TestGetBalanceHandler(t *testing.T) {
	mockClient := &MockETHClient{
		Balance: big.NewInt(1000000000000000000), // 1 ETH
	}
	server := &EVMServer{client: mockClient}

	resp, err := server.GetBalanceHandler(context.Background(), GetBalanceArgs{Address: "0x0000000000000000000000000000000000000000"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data := resp.(map[string]interface{})
	if data["balance"] != "1" {
		t.Errorf("Expected balance 1, got %s", data["balance"])
	}
}

func TestGetTokenBalanceHandler(t *testing.T) {
	mockClient := &MockETHClient{
		CallResult: common.LeftPadBytes(big.NewInt(500).Bytes(), 32),
	}
	server := &EVMServer{client: mockClient}

	resp, err := server.GetTokenBalanceHandler(context.Background(), GetTokenBalanceArgs{
		TokenAddress: "0x0000000000000000000000000000000000000001",
		OwnerAddress: "0x0000000000000000000000000000000000000002",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data := resp.(map[string]interface{})
	if data["balance"] != "500" {
		t.Errorf("Expected balance 500, got %s", data["balance"])
	}
}
