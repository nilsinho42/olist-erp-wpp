package model

type CustomerOrderEnpoint struct {
	CompanyOrderEndpoint
	NomeVendedor string `json:"nomeVendedor"`
}

type Order struct {
	ID            int                  `json:"id"`
	Situacao      int                  `json:"situacao"`
	NumeroPedido  int                  `json:"numeroPedido"`
	DataCriacao   string               `json:"dataCriacao"`
	DataPrevista  string               `json:"dataPrevista"`
	Cliente       CustomerOrderEnpoint `json:"cliente"`
	Valor         string               `json:"valor"`
	Vendedor      Vendedor             `json:"vendedor"`
	Transportador Transportador        `json:"transportador"`
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

type CompanyOrderEndpoint struct {
	TipoCadastro string   `json:"tipoCadastro"`
	ID           int      `json:"id"`
	Codigo       Code     `json:"codigo"`
	TipoPessoa   string   `json:"tipoPessoa"`
	RazaoSocial  string   `json:"nome"`
	NomeFantasia string   `json:"fantasia"`
	CNPJCPF      CNPJCPF  `json:"cpfcnpj"`
	Endereco     Endereco `json:"endereco"`
	Email        Email    `json:"email"`
	Telefone     Telefone `json:"telefone"`
}
