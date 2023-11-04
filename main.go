package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()

	c1 := make(chan *addressResponse)
	c2 := make(chan *addressResponse)

	go func() {
		a1, err := GetAddressAPIByZipcode("22735140", ctx)
		if err != nil {
			panic(err)
		}
		c1 <- a1
	}()

	go func() {
		a2, err := GetAddressCDNByZipcode("22735140", ctx)
		if err != nil {
			panic(err)
		}
		c2 <- a2
	}()

	select {
	case address1 := <-c1:
		fmt.Printf("receive from %s, values: %v", address1.url, address1.address)
	case address2 := <-c2:
		fmt.Printf("receive from %s, values: %v", address2.url, address2.address)
	case <-time.After(time.Second * 1):
		println("timeout")
	}
}

type addressResponse struct {
	address any
	url     string
}

type zipCodeAPIResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func GetAddressAPIByZipcode(zipcode string, ctx context.Context) (*addressResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*1000)
	defer cancel()

	url := "https://viacep.com.br/ws/" + zipcode + "/json/"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var z zipCodeAPIResponse
	err = json.Unmarshal(body, &z)
	if err != nil {
		return nil, err
	}

	return &addressResponse{
		address: &z,
		url:     url,
	}, nil
}

type zipCodeCDNResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
}

func GetAddressCDNByZipcode(zipcode string, ctx context.Context) (*addressResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*1000)
	defer cancel()

	url := "https://opencep.com/v1/" + zipcode + ".json"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var z zipCodeCDNResponse
	err = json.Unmarshal(body, &z)
	if err != nil {
		return nil, err
	}

	return &addressResponse{
		address: &z,
		url:     url,
	}, nil
}
