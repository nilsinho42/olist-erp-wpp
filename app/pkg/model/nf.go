package model

type NF struct {
	ID                 int                  `json:"id"`
	Situacao           string               `json:"situacao"`
	Numero             string               `json:"numero"`
	Serie              string               `json:"serie"`
	ChaveAcesso        string               `json:"chaveAcesso"`
	DataEmissao        string               `json:"dataEmissao"`
	Cliente            CustomerOrderEnpoint `json:"cliente"`
	Valor              float64              `json:"valor"`
	ValorProdutos      float64              `json:"valorProdutos"`
	CodigoRastreamento string               `json:"codigoRastreamento"`
	UrlRastreamento    string               `json:"urlRastreamento"`
	FretePorConta      string               `json:"fretePorConta"`
}

type NFResponse struct {
	Itens     []NF      `json:"itens"`
	Paginacao Paginacao `json:"paginacao"`
}

/*
{"situacao":"6",
"tipo":"S",
"numero":"008429",
"serie":"1",
"chaveAcesso":"35221008804860000190550010000084291892480092",
"dataEmissao":"2022-10-19",
"cliente":
  {"nome":"Casa dos Aquecedores & Cia Ltda",
  "codigo":"",
  "fantasia":"",
  "tipoPessoa":"S",
  "cpfCnpj":"47.545.205\/0001-16",
  "inscricaoEstadual":"636509317116",
  "rg":"",
  "telefone":"(11) 4221-4229",
  "celular":"",
  "email":"",
  "endereco":
    {"endereco":"Rua Alegre",
    "numero":"819",
    "complemento":"",
    "bairro":"Santa Paula",
    "municipio":"São Caetano do Sul",
    "cep":"09.550-250",
    "uf":"SP",
    "pais":"Brasil"},
  "id":753145944},
"enderecoEntrega":null,
"valor":2185,
"valorProdutos":2300,
"valorFrete":0,
"vendedor":null,
"idFormaEnvio":882692566,
"idFormaFrete":0,
"codigoRastreamento":"",
"urlRastreamento":"",
"fretePorConta":"R",
"qtdVolumes":4,
"pesoBruto":22,
"pesoLiquido":0,
"id":893185285,
"ecommerce":
  {"id":0,
  "nome":"",
  "numeroPedidoEcommerce":"",
  "numeroPedidoCanalVenda":"",
  "canalVenda":""}
}
*/
