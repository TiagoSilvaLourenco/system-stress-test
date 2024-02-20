package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado.")
	requests := flag.Int("requests", 0, "Número total de requests.")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas.")
	maxRetries := flag.Int("max-retries", 3, "Número máximo de tentativas para solicitações falhadas.")
	flag.Parse()

	if *url == "" || *requests <= 0 || *concurrency <= 0 || *maxRetries <= 0 {
		fmt.Println("Por favor, forneça valores válidos para URL, requests, concurrency e max-retries.")
		return
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	requestCounter := 0
	successfulRequests := 0
	totalRequests := 0
	statusCodes := make(map[int]int)

	startTime := time.Now()

	doRequest := func() {
		defer wg.Done()

		var err error
		var response *http.Response

		// Realiza tentativas até atingir o limite
		for retries := 0; retries < *maxRetries; retries++ {
			response, err = http.Get(*url)
			if err == nil && response.StatusCode == http.StatusOK {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		mu.Lock()
		defer mu.Unlock()

		if err != nil {
			fmt.Println("Erro na request:", err)
			return
		}
		defer response.Body.Close()

		statusCodes[response.StatusCode]++
		if response.StatusCode == http.StatusOK {
			successfulRequests++
		}
		requestCounter++
		totalRequests++
	}

	// Inicializa as goroutines
	wg.Add(*requests)
	for i := 0; i < *requests; i++ {
		go doRequest()
	}

	// Aguarda a conclusão de todas as goroutines
	wg.Wait()

	// Verifica se é necessário realizar mais solicitações recursivamente
	for successfulRequests < *requests {
		wg.Add(*requests - successfulRequests)
		for i := 0; i < (*requests - successfulRequests); i++ {
			go doRequest()
		}
		wg.Wait()
	}

	// Exibindo o tempo total gasto nos testes
	elapsedTime := time.Since(startTime)

	// Exibindo o relatório principal
	fmt.Printf("Relatório de Teste:\n")
	fmt.Printf("Tempo total gasto: %v\n", elapsedTime)
	fmt.Printf("Quantidade total de requests: %d\n", totalRequests)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", successfulRequests)
	fmt.Printf("Distribuição de códigos de status HTTP:\n")
	for code, count := range statusCodes {
		if code != http.StatusOK {
			fmt.Printf("%d: %d\n", code, count)
			continue
		}
	}
}
