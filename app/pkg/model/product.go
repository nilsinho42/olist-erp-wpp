package model

type Preco struct {
	Preco            float64 `json:"preco"`
	PrecoPromocional float64 `json:"precoPromocional"`
}

type Dimensoes struct {
	Largura           float64 `json:"largura"`
	Altura            float64 `json:"altura"`
	Comprimento       float64 `json:"comprimento"`
	Diametro          float64 `json:"diametro"`
	PesoLiquido       float64 `json:"pesoLiquido"`
	PesoBruto         float64 `json:"pesoBruto"`
	QuantidadeVolumes int     `json:"quantidadeVolumes"`
}

type Product struct {
	ID          int       `json:"id"`
	NomeProduto string    `json:"descricao"`
	Precos      Preco     `json:"precos"`
	Unidade     string    `json:"unidade"`
	Dimensoes   Dimensoes `json:"dimensoes"`
}

type ProductResponse struct {
	Itens     []Product `json:"itens"`
	Paginacao Paginacao `json:"paginacao"`
}

/*
{
  "id": 0,
  "sku": "string",
  "descricao": "string",
  "descricaoComplementar": "string",
  "tipo": "K",
  "situacao": "A",
  "produtoPai": {
    "id": 0,
    "sku": "string",
    "descricao": "string"
  },
  "unidade": "string",
  "unidadePorCaixa": "string",
  "ncm": "string",
  "gtin": "string",
  "origem": "0",
  "garantia": "string",
  "observacoes": "string",
  "categoria": {
    "id": 0,
    "nome": "string",
    "caminhoCompleto": "string"
  },
  "marca": {
    "id": 0,
    "nome": "string"
  },
  "dimensoes": {
    "embalagem": {
      "id": 0,
      "tipo": 0,
      "descricao": "string"
    },
    "largura": 0,
    "altura": 0,
    "comprimento": 0,
    "diametro": 0,
    "pesoLiquido": 0,
    "pesoBruto": 0,
    "quantidadeVolumes": 0
  },
  "precos": {
    "preco": 0,
    "precoPromocional": 0,
    "precoCusto": 0,
    "precoCustoMedio": 0
  },
  "estoque": {
    "controlar": true,
    "sobEncomenda": true,
    "diasPreparacao": 0,
    "localizacao": "string",
    "minimo": 0,
    "maximo": 0,
    "quantidade": 0
  },
  "fornecedores": [
    {
      "id": 0,
      "nome": "string",
      "codigoProdutoNoFornecedor": "string"
    }
  ],
  "seo": {
    "titulo": "string",
    "descricao": "string",
    "keywords": [
      "string"
    ],
    "linkVideo": "string",
    "slug": "string"
  },
  "tributacao": {
    "gtinEmbalagem": "string",
    "valorIPIFixo": 0,
    "classeIPI": "string"
  },
  "anexos": [
    {
      "url": "string",
      "externo": true
    }
  ],
  "variacoes": [
    {
      "id": 0,
      "descricao": "string",
      "sku": "string",
      "gtin": "string",
      "precos": {
        "preco": 0,
        "precoPromocional": 0,
        "precoCusto": 0,
        "precoCustoMedio": 0
      },
      "estoque": {
        "controlar": true,
        "sobEncomenda": true,
        "diasPreparacao": 0,
        "localizacao": "string",
        "minimo": 0,
        "maximo": 0,
        "quantidade": 0
      },
      "grade": [
        {
          "chave": "string",
          "valor": "string"
        }
      ]
    }
  ],
  "kit": [
    {
      "produto": {
        "id": 0,
        "sku": "string",
        "descricao": "string"
      },
      "quantidade": 0
    }
  ],
  "producao": {
    "produtos": [
      {
        "produto": {
          "id": 0,
          "sku": "string",
          "descricao": "string"
        },
        "quantidade": 0
      }
    ],
    "etapas": [
      "string"
    ]
  }
}

{
      "id": 0,
      "sku": "string",
      "descricao": "string",
      "tipo": "K",
      "situacao": "A",
      "dataCriacao": "string",
      "dataAlteracao": "string",
      "unidade": "string",
      "gtin": "string",
      "precos": {
        "preco": 0,
        "precoPromocional": 0,
        "precoCusto": 0,
        "precoCustoMedio": 0
      }


type CustomerOrderEnpoint struct {
	CompanyOrderEndpoint
	NomeVendedor string `json:"nomeVendedor"`
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

{"id":905160935,
"sku":"CSEP",
"descricao":"Cascata Inox Sertech Easy Profissional ",
"descricaoComplementar":"",
"tipo":"S",
"situacao":"A",
"produtoPai":null,
"unidade":"PÇ",
"unidadePorCaixa":"",
"ncm":"7326.90.90",
"gtin":"",
"origem":"0",
"garantia":"",
"observacoes":"",
"categoria":null,
"marca":null,
"dimensoes":
  {"embalagem":
    {"id":null,
    "tipo":2,
    "descricao":""},
"largura":55,
"altura":40,
"comprimento":105,
"diametro":0,
"pesoLiquido":12,
"pesoBruto":12,
"quantidadeVolumes":1},
"precos":
  {"preco":830,
  "precoPromocional":0,
  "precoCusto":0,
  "precoCustoMedio":0},
"estoque":
  {"controlar":false,
  "sobEncomenda":false,
  "diasPreparacao":0,
  "localizacao":"",
  "minimo":0,
  "maximo":0,
  "quantidade":0},
"fornecedores":[],
"seo":
  {"titulo":"",
  "descricao":"",
  "keywords":[""],
  "linkVideo":"",
  "slug":""},
"tributacao":
  {"gtinEmbalagem":"",
  "valorIPIFixo":0,
  "classeIPI":""},
"anexos":[],
"variacoes":[],
"kit":[],
"producao":null
}%

Response:
{"id":915587818,
"sku":"CSEG",
"descricao":"Cascata Inox Sertech Easy - Aço 316",
"descricaoComplementar":"",
"tipo":"S","situacao":"A","produtoPai":null,"unidade":"PÇ","unidadePorCaixa":"","ncm":"7326.90.90","gtin":"","origem":"0","garantia":"","observacoes":"","categoria":null,"marca":null,"dimensoes":{"embalagem":{"id":null,"tipo":2,"descricao":""},"largura":0,"altura":0,"comprimento":0,"diametro":0,"pesoLiquido":12,"pesoBruto":12,"quantidadeVolumes":1},"precos":{"preco":1300,"precoPromocional":0,"precoCusto":0,"precoCustoMedio":0},"estoque":{"controlar":true,"sobEncomenda":false,"diasPreparacao":0,"localizacao":"","minimo":0,"maximo":0,"quantidade":0},"fornecedores":[],"seo":{"titulo":"","descricao":"","keywords":[""],"linkVideo":"","slug":""},"tributacao":{"gtinEmbalagem":"","valorIPIFixo":0,"classeIPI":""},"anexos":[],"variacoes":[],"kit":[],"producao":null}
*/
