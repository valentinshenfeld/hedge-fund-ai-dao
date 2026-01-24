package main

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
    // "github.com/hedge-fund/internal/mcp"
    // "github.com/hedge-fund/internal/a2a"
)

func main() {
	ctx := context.Background()

    // 1. Инициализация MCP клиента для X (Twitter)
    twitterMCP := mcp.NewClient("http://mcp-server-x:8080")
    
    // 2. Инициализация Google GenAI (Gemini)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err!= nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro")
    
    // 3. Регистрация инструментов (Tools) из MCP в модель
    model.Tools =*genai.Tool{
        twitterMCP.AsGeminiTool(), // Авто-конвертация MCP определений в Gemini Tools
    }

	// 4. Запуск A2A сервера для приема задач
    agentServer := a2a.NewServer(":50051", "AnalystAgent")
    agentServer.OnTask("ANALYZE_TOKEN", func(payloadbyte) (byte, error) {
        
        // Логика рассуждения
        prompt := "Analyze the sentiment for token " + string(payload) + " using the Twitter tool."
        resp, err := model.GenerateContent(ctx, genai.Text(prompt))
        if err!= nil {
            return nil, err
        }
        
        // Обработка Function Calls (если модель решила вызвать Twitter)
        //... (код обработки вызовов инструментов)
        
        returnbyte(resp.Candidates.Content.Parts.(genai.Text)), nil
    })

    log.Println("Analyst Agent running...")
    agentServer.Start()
}