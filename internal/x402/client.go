package x402

import (
	"context"
	"net/http"
	"github.com/ethereum/go-ethereum/crypto"
    //... imports
)

type PaymentClient struct {
	PrivateKey *ecdsa.PrivateKey
	Client     *http.Client
}

// DoWithPayment выполняет запрос и обрабатывает 402 ответ автоматически
func (pc *PaymentClient) DoWithPayment(req *http.Request) (*http.Response, error) {
	// 1. Первый запрос
	resp, err := pc.Client.Do(req)
	if err!= nil {
		return nil, err
	}

	// 2. Если статус не 402, возвращаем как есть
	if resp.StatusCode!= 402 {
		return resp, nil
	}
	defer resp.Body.Close()

	// 3. Парсинг заголовков x402 (WWW-Authenticate)
	paymentDetails, err := parseX402Headers(resp.Header)
	if err!= nil {
		return nil, err
	}

	// 4. Подписание платежа (EIP-712 или простая подпись транзакции)
	proof, err := pc.signPayment(paymentDetails)
	if err!= nil {
		return nil, err
	}

	// 5. Повтор запроса с заголовком Authorization
	req.Header.Set("Authorization", "X402 " + proof)
    // Необходимо "перемотать" body запроса, если он был
    
	return pc.Client.Do(req)
}

func (pc *PaymentClient) signPayment(details PaymentDetails) (string, error) {
    // Логика подписи транзакции или создания ваучера
	return "signed_proof_hex", nil
}