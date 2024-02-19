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

	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		fmt.Println("Por favor, forneça valores válidos para URL, requests e concurrency.")
		return
	}

	var wg sync.WaitGroup
	requestCounter := 0
	successfulRequests := 0
	statusCodes := make(map[int]int)
	errorChannel := make(chan error, *requests)

	startTime := time.Now()

	doRequest := func() {
		defer wg.Done()

		response, err := http.Get(*url)
		if err != nil {
			errorChannel <- err
			return
		}
		defer response.Body.Close()

		statusCodes[response.StatusCode]++

		if response.StatusCode == http.StatusOK {
			successfulRequests++
		}

		requestCounter++
	}

	bucket := make(chan struct{}, *concurrency)

	for i := 0; i < *requests; i++ {
		wg.Add(1)
		go func() {
			bucket <- struct{}{}
			doRequest()
			<-bucket
		}()
	}

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	for err := range errorChannel {
		fmt.Println(err)
	}

	elapsedTime := time.Since(startTime)

	fmt.Printf("Relatório de Teste:\n")
	fmt.Printf("Tempo total gasto: %v\n", elapsedTime)
	fmt.Printf("Quantidade total de requests: %d\n", requestCounter)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", successfulRequests)
	fmt.Printf("Distribuição de códigos de status HTTP:\n")
	for code, count := range statusCodes {
		fmt.Printf("%d: %d\n", code, count)
	}
}
