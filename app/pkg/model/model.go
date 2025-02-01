package model

import (
	"fmt"
	"strings"
)

type CNPJCPF string
type Email string
type Telefone string
type Code string
type Data string

type Endereco struct {
	Rua         string `json:"endereco"`
	Numero      string `json:"numero"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Municipio   string `json:"municipio"`
	Cep         string `json:"cep"`
	Uf          string `json:"uf"`
	Pais        string `json:"pais"`
}

// Supplier struct with custom types
type Product struct {
	DataCriacao      string  `json:"dataCriacao"`
	NomeProduto      string  `json:"nomeProduto"`
	Codigo           Code    `json:"codigo"`
	Preco            float64 `json:"preco"`
	PrecoPromocional float64 `json:"precoPromocional"`
	Unidade          string  `json:"unidade"`
	GTIN             string  `json:"gtin"`
}

type Vendedor struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

type FormaEnvio struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

type Transportador struct {
	ID                 int        `json:"id"`
	Nome               string     `json:"nome"`
	FretePorConta      string     `json:"fretePorConta"`
	FormaEnvio         FormaEnvio `json:"formaEnvio"`
	FormaFrete         string     `json:"formaFrete"`
	CodigoRastreamento string     `json:"codigoRastreamento"`
	UrlRastreamento    string     `json:"urlRastreamento"`
}

type Paginacao struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

func ValidateCPFCNPJ(input string) (string, error) {
	// If CNPJ: 15.049.188/0001-30 => 15.049.188/0001-30, nil
	// If CNPJ: 15049188000130 => 15.049.188/0001-30, nil
	// If CPF: 123.456.789-09 => 123.456.789-09, nil
	// If CPF: 12345678909 => 123.456.789-09, nil
	// any other format:  given, error

	// If no ., / or - AND length is 14, split in 3 parts and add . and / or - in the right places
	// If no . or - AND length is 11, split in 4 parts and add . and - in the right places

	input = strings.ReplaceAll(input, ".", "")
	input = strings.ReplaceAll(input, "-", "")
	input = strings.ReplaceAll(input, "/", "")
	if len(input) == 14 {
		fmt.Printf("%s.%s.%s/%s-%s\n", input[:2], input[2:5], input[5:8], input[8:12], input[12:14])
		return fmt.Sprintf("%s.%s.%s/%s-%s", input[:2], input[2:5], input[5:8], input[8:12], input[12:14]), nil
	} else if len(input) == 11 {
		return fmt.Sprintf("%s.%s.%s-%s", input[:3], input[3:6], input[6:9], input[9:11]), nil
	}
	return input, fmt.Errorf("invalid CPF/CNPJ format")

}
