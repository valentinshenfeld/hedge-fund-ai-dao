package main

import (
	"context"
	"log"
	"os"
	"google.golang.org/adk/a2a"
	"google.golang.org/adk/mcp"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Model interface {
	GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
}

type AnalystAgent struct {
	model Model
	ctx   context.Context
}

func (a *AnalystAgent) AnalyzeTokenHandler(payload []byte) ([]byte, error) {
	// Логика рассуждения
	prompt := "Analyze the sentiment for token " + string(payload) + " using the Twitter tool."
	resp, err := a.model.GenerateContent(a.ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	// Обработка Function Calls (если модель решила вызвать Twitter)
	// В реальной системе здесь будет цикл обработки вызовов инструментов
	
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				return []byte(text), nil
			}
		}
	}

	return []byte("No analysis result"), nil
}

func SetupAnalystAgent(ctx context.Context, apiKey string) (*AnalystAgent, error) {
	// 1. Инициализация MCP клиента для X (Twitter)
	twitterMCP := mcp.NewClient("http://mcp-server-x:8080")

	// 2. Инициализация Google GenAI (Gemini)
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-1.5-pro")

	// 3. Регистрация инструментов (Tools) из MCP в модель
	model.Tools = []*genai.Tool{
		twitterMCP.AsGeminiTool().(*genai.Tool),
	}

	return &AnalystAgent{
		model: model,
		ctx:   ctx,
	}, nil
}

func main() {
	ctx := context.Background()
	agent, err := SetupAnalystAgent(ctx, os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// 4. Запуск A2A сервера для приема задач
	agentServer := a2a.NewServer(":50051", "AnalystAgent")
	agentServer.OnTask("ANALYZE_TOKEN", agent.AnalyzeTokenHandler)

	log.Println("Analyst Agent running...")
	agentServer.Start()
}