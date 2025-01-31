package model

type AccountReceivable struct {
	Cliente    CustomerOrderEnpoint `json:"cliente"`
	Doc        string               `json:"numeroDocumento"`
	Vencimento Data                 `json:"dataVencimento"`
	Situacao   string               `json:"situacao"`
	Valor      float64              `json:"valor"`
}

type AccountReceivableResponse struct {
	Itens     []AccountReceivable `json:"itens"`
	Paginacao Paginacao           `json:"paginacao"`
}

type AccountPayable struct {
}

type AccountPayableResponse struct {
	Itens     []AccountPayable `json:"itens"`
	Paginacao Paginacao        `json:"paginacao"`
}

/*
{"id":894650003,
"situacao":"pago",
"data":"2024-10-23",
"dataVencimento":"2022-03-24",
"historico":"Ref. a NF nº 8068, MARCUCCI REVESTIMENTOS PARA PISCINAS LTDA",
"valor":742,
"numeroDocumento":"008068\/01",
"numeroBanco":null,
"serieDocumento":1,
"cliente":
	{"nome":"MARCUCCI REVESTIMENTOS PARA PISCINAS LTDA",
	"codigo":"",
	"fantasia":"",
	"tipoPessoa":"J",
	"cpfCnpj":"52.454.212\/0001-42",
	"inscricaoEstadual":"587253090110",
	"rg":"",
	"telefone":"(19) 3523-1125",
	"celular":"",
	"email":"",
	"endereco":
		{"endereco":"NOVE",
		"numero":"671",
		"complemento":"",
		"bairro":"CENTRO",
		"municipio":"Rio Claro",
		"cep":"13.500-080",
		"uf":"SP",
		"pais":"Brasil"},
	"id":753146269}}
*/
