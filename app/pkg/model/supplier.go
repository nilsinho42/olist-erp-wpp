package model

// Using composition to create Supplier and Customer
type Supplier struct {
	Company
}

type SupplierResponse struct {
	Itens []struct {
		Nome              string `json:"nome"`
		Codigo            string `json:"codigo"`
		Fantasia          string `json:"fantasia"`
		TipoPessoa        string `json:"tipoPessoa"`
		CpfCnpj           string `json:"cpfCnpj"`
		InscricaoEstadual string `json:"inscricaoEstadual"`
		Rg                string `json:"rg"`
		Telefone          string `json:"telefone"`
		Celular           string `json:"celular"`
		Email             string `json:"email"`
		Endereco          struct {
			Endereco    string `json:"endereco"`
			Numero      string `json:"numero"`
			Complemento string `json:"complemento"`
			Bairro      string `json:"bairro"`
			Municipio   string `json:"municipio"`
			Cep         string `json:"cep"`
			Uf          string `json:"uf"`
			Pais        string `json:"pais"`
		} `json:"endereco"`
		ID       int `json:"id"`
		Vendedor struct {
			ID   int    `json:"id"`
			Nome string `json:"nome"`
		} `json:"vendedor"`
		Situacao        string `json:"situacao"`
		DataCriacao     string `json:"dataCriacao"`
		DataAtualizacao string `json:"dataAtualizacao"`
	} `json:"itens"`
	Paginacao struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	} `json:"paginacao"`
}
