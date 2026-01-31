package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/mcp"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

// Configuration for Agent URLs (In K8s these are DNS names)
var agents = map[string]string{
	"analyst":          "http://analyst:50051",
	"agent-risk":       "http://risk-analyst:50052",
	"agent-trader":     "http://trader:50053",
	"mcp-server-x":     "http://mcp-server-x:50054",
	"mcp-server-evm":   "http://mcp-server-evm:50055",
}

type AgentRequest struct {
	Message string `json:"message"`
}

type AgentResponse struct {
	Response string `json:"response"`
	Status   string `json:"status"`
}

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type HTTPClient interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type WorkflowManager struct {
	rdb    RedisClient
	ctx    context.Context
	client HTTPClient
}

func main() {
	// 1. Init Redis (Short-lived persistent memory)
	rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"), // e.g., "redis:6379"
	})

	manager := &WorkflowManager{
		rdb:    rdb,
		ctx:    context.Background(),
		client: http.DefaultClient,
	}

	http.HandleFunc("/start_cycle", manager.HandleBlogCycle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Orchestrator listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (m *WorkflowManager) HandleBlogCycle(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		http.Error(w, "Topic is required", http.StatusBadRequest)
		return
	}

	go m.RunWorkflow(topic) // Async execution

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Workflow started for topic: %s", topic)
}

// The Core Logic: Agent-to-Agent (A2A) Coordination
func (m *WorkflowManager) RunWorkflow(topic string) {
	log.Printf("Starting workflow for: %s", topic)

	// Step 1: Analyst
	researchData := m.CallAgent("analyst", "Research this topic deeply: "+topic)
	m.rdb.Set(m.ctx, topic+":research", researchData, 24*time.Hour)

	// Step 2: Risk Agent
	riskData := m.CallAgent("agent-risk", "Analyze the risk of: "+researchData)
	m.rdb.Set(m.ctx, topic+":risk", riskData, 24*time.Hour)

	// Step 3: Trader Agent
	tradeData := m.CallAgent("agent-trader", "Execute trade based on research: "+researchData)
	m.rdb.Set(m.ctx, topic+":trade", tradeData, 24*time.Hour)

	// Step 4: MCP Server X
	xData := m.CallAgent("mcp-server-x", "Post summary to X: "+researchData)
	m.rdb.Set(m.ctx, topic+":x", xData, 24*time.Hour)

	// Step 5: MCP Server EVM
	evmData := m.CallAgent("mcp-server-evm", "Check EVM status: "+researchData)
	m.rdb.Set(m.ctx, topic+":evm", evmData, 24*time.Hour)
}

func (m *WorkflowManager) CallAgent(agentName string, prompt string) string {
	url, ok := agents[agentName]
	if !ok {
		log.Printf("Unknown agent: %s", agentName)
		return ""
	}

	reqBody, _ := json.Marshal(AgentRequest{Message: prompt})

	resp, err := m.client.Post(url+"/task", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error calling %s: %v", agentName, err)
		return ""
	}
	defer resp.Body.Close()

	var res AgentResponse
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Response
}