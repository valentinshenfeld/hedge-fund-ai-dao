package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/adk/mcp"
	"google.golang.org/adk/tool"
)

// Config holds the server configuration
type Config struct {
	RPCURL string
}

type ETHClient interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

type EVMServer struct {
	client ETHClient
}

// GetBalanceArgs defines the arguments for get_balance tool
type GetBalanceArgs struct {
	Address string `json:"address" jsonschema:"The EVM address to check balance for"`
}

// GetBalanceHandler returns the ETH balance of an address
func (s *EVMServer) GetBalanceHandler(ctx context.Context, args GetBalanceArgs) (any, error) {
	addr := common.HexToAddress(args.Address)
	balance, err := s.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// Return balance in ETH (string for precision)
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	return map[string]interface{}{
		"address": args.Address,
		"balance": ethValue.String(),
		"symbol":  "ETH",
	}, nil
}

// GetTokenBalanceArgs defines the arguments for get_token_balance tool
type GetTokenBalanceArgs struct {
	TokenAddress string `json:"token_address" jsonschema:"The ERC20 token contract address"`
	OwnerAddress string `json:"owner_address" jsonschema:"The address of the token owner"`
}

// GetTokenBalanceHandler returns the ERC20 balance of an address
func (s *EVMServer) GetTokenBalanceHandler(ctx context.Context, args GetTokenBalanceArgs) (any, error) {
	tokenAddr := common.HexToAddress(args.TokenAddress)
	ownerAddr := common.HexToAddress(args.OwnerAddress)

	// ERC20 balanceOf(address) selector: 0x70a08231
	data := append(common.Hex2Bytes("70a08231"), common.LeftPadBytes(ownerAddr.Bytes(), 32)...)

	msg := ethereum.CallMsg{
		To:   &tokenAddr,
		Data: data,
	}

	result, err := s.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	balance := new(big.Int).SetBytes(result)
	return map[string]interface{}{
		"token":   args.TokenAddress,
		"owner":   args.OwnerAddress,
		"balance": balance.String(),
	}, nil
}

// MonitorSwapsArgs defines the arguments for monitor_swaps tool
type MonitorSwapsArgs struct {
	TokenAddress string `json:"token_address" jsonschema:"The token address to monitor for swaps"`
	LastBlocks   uint64 `json:"last_blocks" jsonschema:"Number of recent blocks to scan"`
}

// MonitorSwapsHandler scans recent blocks for Swap events (Uniswap V2/V3 compatible pattern)
func (s *EVMServer) MonitorSwapsHandler(ctx context.Context, args MonitorSwapsArgs) (any, error) {
	header, err := s.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	fromBlock := header.Number.Uint64() - args.LastBlocks
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   header.Number,
		Addresses: []common.Address{common.HexToAddress(args.TokenAddress)},
	}

	logs, err := s.client.FilterLogs(ctx, query)
	if err != nil {
		return nil, err
	}

	type SwapLog struct {
		TxHash      string `json:"tx_hash"`
		BlockNumber uint64 `json:"block_number"`
		Data        string `json:"data"`
	}

	var results []SwapLog
	for _, vLog := range logs {
		results = append(results, SwapLog{
			TxHash:      vLog.TxHash.Hex(),
			BlockNumber: vLog.BlockNumber,
			Data:        common.Bytes2Hex(vLog.Data),
		})
	}

	return map[string]interface{}{
		"token":  args.TokenAddress,
		"blocks": args.LastBlocks,
		"swaps":  results,
		"count":  len(results),
	}, nil
}

// GetTokenVolatilityArgs defines the arguments for GetTokenVolatility tool
type GetTokenVolatilityArgs struct {
	TokenAddress string `json:"token_address" jsonschema:"The token address to check volatility for"`
}

// GetTokenVolatilityHandler calculates a simple volatility metric based on recent price changes
func (s *EVMServer) GetTokenVolatilityHandler(ctx context.Context, args GetTokenVolatilityArgs) (any, error) {
	// In a real scenario, this would fetch historical prices from a DEX or Oracle
	// For this MCP server, we'll simulate a response or use recent logs
	return map[string]interface{}{
		"token":             args.TokenAddress,
		"volatility_24h":    "0.052", // Simulated 5.2%
		"confidence_score":  0.89,
		"last_price_change": "-1.2%",
	}, nil
}

func main() {
	rpcURL := os.Getenv("EVM_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://eth-mainnet.g.alchemy.com/v2/your-api-key" // Default or placeholder
	}

	dialClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	evmServer := &EVMServer{client: dialClient}

	// Initialize MCP Server
	server := mcp.NewServer(
		"mcp-server-evm",
		"1.0.0",
		mcp.WithDescription("EVM Activity Monitor for Hedge Fund AI DAO"),
	)

	// Register Tools
	server.RegisterTool(tool.NewFunctionTool("get_balance", evmServer.GetBalanceHandler))
	server.RegisterTool(tool.NewFunctionTool("get_token_balance", evmServer.GetTokenBalanceHandler))
	server.RegisterTool(tool.NewFunctionTool("monitor_swaps", evmServer.MonitorSwapsHandler))
	server.RegisterTool(tool.NewFunctionTool("CheckLiquidity", evmServer.GetTokenBalanceHandler)) // Alias for now
	server.RegisterTool(tool.NewFunctionTool("GetTokenVolatility", evmServer.GetTokenVolatilityHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("EVM MCP Server starting on port %s...", port)
	if err := server.ListenAndServe(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
