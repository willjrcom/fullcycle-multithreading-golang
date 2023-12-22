package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const timeout = 1 * time.Second

type AddressBrasilAPI struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"street"`
	Bairro     string `json:"neighborhood"`
	Cidade     string `json:"city"`
	UF         string `json:"state"`
}

type AddressViacep struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Cidade     string `json:"localidade"`
	UF         string `json:"uf"`
}

func main() {
	cep := "01153000"

	// Use WaitGroup para esperar que ambas as goroutines concluam
	var wg sync.WaitGroup
	wg.Add(2)

	// Canal para receber o resultado da API 1
	ch1 := make(chan *AddressBrasilAPI, 1)

	// Goroutine para a primeira API
	go func() {
		defer wg.Done()
		result, err := fetchFromBrasilAPI(cep)
		if err != nil {
			fmt.Println("Erro na API 1:", err)
			return
		}
		ch1 <- result
	}()

	// Canal para receber o resultado da API 2
	ch2 := make(chan *AddressViacep, 1)

	// Goroutine para a segunda API
	go func() {
		defer wg.Done()
		result, err := fetchFromViacep(cep)
		if err != nil {
			fmt.Println("Erro na API 2:", err)
			return
		}
		ch2 <- result
	}()

	// Use Select para esperar o resultado de qualquer API
	select {
	case result := <-ch1:
		fmt.Println("Resultado da API 1:")
		displayResultBrasilAPI(result, "API 1")
	case result := <-ch2:
		fmt.Println("Resultado da API 2:")
		displayResultViacep(result, "API 2")
	case <-time.After(timeout):
		fmt.Println("Timeout: nenhuma resposta dentro do tempo limite")
	}

	// Espera ambas as goroutines concluírem
	wg.Wait()
}

func fetchFromBrasilAPI(cep string) (*AddressBrasilAPI, error) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep

	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result AddressBrasilAPI
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func fetchFromViacep(cep string) (*AddressViacep, error) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"

	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))
	var result AddressViacep
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func displayResultBrasilAPI(result *AddressBrasilAPI, apiName string) {
	fmt.Printf("API: %s\n", apiName)
	fmt.Printf("CEP: %s\n", result.Cep)
	fmt.Printf("Logradouro: %s\n", result.Logradouro)
	fmt.Printf("Bairro: %s\n", result.Bairro)
	fmt.Printf("Cidade: %s\n", result.Cidade)
	fmt.Printf("Estado: %s\n", result.UF)
}

func displayResultViacep(result *AddressViacep, apiName string) {
	fmt.Printf("API: %s\n", apiName)
	fmt.Printf("CEP: %s\n", result.Cep)
	fmt.Printf("Logradouro: %s\n", result.Logradouro)
	fmt.Printf("Bairro: %s\n", result.Bairro)
	fmt.Printf("Cidade: %s\n", result.Cidade)
	fmt.Printf("Estado: %s\n", result.UF)
}
