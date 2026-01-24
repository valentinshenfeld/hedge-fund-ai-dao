package main

import (
    "github.com/smartcontractkit/cre-sdk/go/workflow"
    "github.com/smartcontractkit/cre-sdk/go/capabilities/evm"
    "encoding/json"
)

// Структура входных данных от агента
type TradeSignal struct {
    TokenAddress string `json:"token"`
    Amount       string `json:"amount"`
    TargetPrice  string `json:"target_price"`
}

func main() {
    // Определение воркфлоу
    app := workflow.New()

    app.Trigger("http_trigger", func(ctx workflow.Context, event workflow.Event) error {
        // 1. Парсинг сигнала от AI Агента
        var signal TradeSignal
        json.Unmarshal(event.Data, &signal)

        // 2. Верификация цены (используем Capability DON)
        // Получаем цену из надежного источника (Chainlink Data Streams)
        price, err := workflow.Cap("fetch-price").Request(map[string]interface{}{
            "asset": signal.TokenAddress,
        }).Await()
        
        if err!= nil {
            return err // Консенсус не достигнут или ошибка сети
        }

        // 3. Сравнение цены AI с рыночной ценой (защита от галлюцинаций)
        if!isPriceValid(signal.TargetPrice, price) {
             return workflow.Error("Price deviation too high")
        }

        // 4. Формирование транзакции
        // Кодирование вызова функции invest() смарт-контракта
        txPayload, _ := evm.Encode("invest(address,uint256)", signal.TokenAddress, signal.Amount)

        // 5. Отправка транзакции через Write Capability
        // Это действие требует консенсуса DON
        _, err = workflow.Cap("execute-onchain").Request(map[string]interface{}{
            "to": "0xAssetManagerAddress...", // Адрес AssetManager
            "data": txPayload,
        }).Await()

        return err
    })

    app.Run()
}

func isPriceValid(agentPrice string, oraclePrice interface{}) bool {
    // Логика проверки отклонения (slippage tolerance)
    return true 
}