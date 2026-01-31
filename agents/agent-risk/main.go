package main

import (
	"context"
	"encoding/json"
	"fmt"
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

type RiskAgent struct {
	model Model
	ctx   context.Context
}

func (a *RiskAgent) ValidateRiskHandler(payload []byte) ([]byte, error) {
	log.Printf("Проверка рисков для стратегии: %s", string(payload))

	// ШАГ 1: Сбор данных через MCP (вызывается автоматически моделью)
	prompt := "Проверь ликвидность и проскальзывание для этой сделки: " + string(payload) +
		". Если проскальзывание > 2% или ликвидность < $100k, отклони."

	resp, err := a.model.GenerateContent(a.ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	// ШАГ 2: Анализ вердикта LLM
	var result struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}

	// Assuming the response part is Text
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				json.Unmarshal([]byte(text), &result)
				break
			}
		}
	}

	if result.Status == "fail" {
		log.Printf("РИСК ОТКЛОНЕН: %s", result.Reason)
		return nil, fmt.Errorf(result.Reason)
	}

	log.Println("РИСК ПРОЙДЕН: Сделка безопасна.")
	return []byte("PASS"), nil
}

func main() {
	ctx := context.Background()

	// 1. Когнитивное ядро (Gemini Flash для быстрой проверки условий)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	model := client.GenerativeModel("gemini-1.5-flash")

	// 2. Инструменты контроля: доступ к блокчейн-метрикам через MCP
	evmMCP := mcp.NewClient("http://mcp-server-evm:8080")
	model.Tools = []*genai.Tool{evmMCP.AsGeminiTool().(*genai.Tool)}

	agent := &RiskAgent{
		model: model,
		ctx:   ctx,
	}

	// 3. Запуск сервера Риск-Менеджера
	riskServer := a2a.NewServer(":50053", "RiskAgent")
	riskServer.OnTask("VALIDATE_RISK", agent.ValidateRiskHandler)

	riskServer.Start()
}