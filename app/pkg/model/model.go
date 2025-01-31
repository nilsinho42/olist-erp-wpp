package model

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

type NF struct {
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
