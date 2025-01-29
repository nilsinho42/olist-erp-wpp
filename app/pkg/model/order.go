package model

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

type Order struct {
	ID            int           `json:"id"`
	Situacao      int           `json:"situacao"`
	NumeroPedido  int           `json:"numeroPedido"`
	Ecommerce     string        `json:"ecommerce"`
	DataCriacao   string        `json:"dataCriacao"`
	DataPrevista  string        `json:"dataPrevista"`
	Cliente       Customer      `json:"cliente"`
	Valor         string        `json:"valor"`
	Vendedor      Vendedor      `json:"vendedor"`
	Transportador Transportador `json:"transportador"`
}

type Paginacao struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type OrderResponse struct {
	Itens     []Order   `json:"itens"`
	Paginacao Paginacao `json:"paginacao"`
}
