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
	flag.Parse()

	// Validando parâmetros
	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		fmt.Println("Por favor, forneça valores válidos para URL, requests e concurrency.")
		return
	}

	// Inicializando variáveis
	var wg sync.WaitGroup
	requestCounter := 0
	successfulRequests := 0
	statusCodes := make(map[int]int)

	// Iniciando o cronômetro
	startTime := time.Now()

	// Função para realizar uma request
	doRequest := func() {
		defer wg.Done()

		response, err := http.Get(*url)
		if err != nil {
			fmt.Println("Erro ao realizar a request:", err)
			return
		}
		defer response.Body.Close()

		// Contabilizando o status code
		statusCodes[response.StatusCode]++

		// Contabilizando a request bem-sucedida
		if response.StatusCode == http.StatusOK {
			successfulRequests++
		}

		// Incrementando o contador total de requests
		requestCounter++
	}

	// Iniciando as chamadas concorrentes
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go doRequest()
	}

	// Aguardando a conclusão de todas as chamadas
	wg.Wait()

	// Calculando o tempo total gasto na execução
	elapsedTime := time.Since(startTime)

	// Exibindo o relatório
	fmt.Printf("Relatório de Teste:\n")
	fmt.Printf("Tempo total gasto: %v\n", elapsedTime)
	fmt.Printf("Quantidade total de requests: %d\n", requestCounter)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", successfulRequests)
	fmt.Printf("Distribuição de códigos de status HTTP:\n")
	for code, count := range statusCodes {
		fmt.Printf("%d: %d\n", code, count)
	}
}
