package model

type CNPJCPF string
type Email string
type Telefone string
type Code string
type Data string

type Endereco struct {
	Rua         string `json:"rua"`
	Numero      string `json:"numero"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Municipio   string `json:"municipio"`
	Cep         string `json:"cep"`
	Uf          string `json:"uf"`
	Pais        string `json:"pais"`
}

// Supplier struct with custom types
type Company struct {
	TipoCadastro string   `json:"tipoCadastro"`
	ID           int      `json:"id"`
	Codigo       Code     `json:"codigo"`
	TipoPessoa   string   `json:"tipoPessoa"`
	RazaoSocial  string   `json:"razaoSocial"`
	NomeFantasia string   `json:"nomeFantasia"`
	CNPJCPF      CNPJCPF  `json:"cnpjCpf"`
	Endereco     Endereco `json:"endereco"`
	Email        Email    `json:"email"`
	Telefone     Telefone `json:"telefone"`
}

type Product struct {
	DataCriacao      string  `json:"dataCriacao"`
	NomeProduto      string  `json:"nomeProduto"`
	Codigo           Code    `json:"codigo"`
	Preco            float64 `json:"preco"`
	PrecoPromocional float64 `json:"precoPromocional"`
	Unidade          string  `json:"unidade"`
	GTIN             string  `json:"gtin"`
}

type Order struct {
	ID              int     `json:"id"`
	NumeroPedido    string  `json:"numeroPedido"`
	NumeroEcommerce string  `json:"numeroEcommerce"`
	DataPedido      Data    `json:"dataPedido"`
	DataPrevista    Data    `json:"dataPrevista"`
	NomeCliente     string  `json:"nomeCliente"`
	Valor           float64 `json:"valor"`
	IDVendedor      int     `json:"idVendedor"`
	NomeVendedor    string  `json:"nomeVendedor"`
	Situacao        string  `json:"situacao"`
	CodigoRastreio  Code    `json:"codigoRastreio"`
}

type NF struct {
}

type AccountReceivable struct {
	NomeCliente string  `json:"nomeCliente"`
	Doc         string  `json:"doc"`
	Detalhes    string  `json:"detalhes"`
	Vencimento  Data    `json:"vencimento"`
	Situacao    string  `json:"situacao"`
	Valor       float64 `json:"valor"`
}

type AccountPayable struct {
}
