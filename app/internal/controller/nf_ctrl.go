package controller

import (
	"app/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type NFController struct {
	nfs                map[string][]*model.NF
	customerController *CustomerController
	token              string
	mu                 sync.RWMutex
}

func NewNFController(token string) *NFController {
	return &NFController{
		nfs:                make(map[string][]*model.NF),
		customerController: NewCustomerController(token),
		token:              token,
	}
}

func (c *NFController) GetByName(ctx context.Context, name string) ([]*model.NF, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if nf, exists := c.nfs[name]; exists {
		return nf, nil
	}
	// Need to convert name to cpfcnpj for search
	cpfcnpj := c.customerController.GetCPFCNPJFromPartialName(ctx, name)
	if cpfcnpj == "" {
		return nil, fmt.Errorf("no customer found with name: %s", name)
	}

	nfs, err := c.GetByCPFCNPJ(ctx, cpfcnpj)
	if err != nil {
		return nil, fmt.Errorf("failed to get NFs by CPFCNPJ: %w", err)
	}
	return nfs, nil
}

func (c *NFController) GetByNumero(ctx context.Context, numero string) ([]*model.NF, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if nf, exists := c.nfs[numero]; exists {
		return nf, nil
	}
	url := fmt.Sprintf("%s/notas?numero=%s", OlistERPURL, numero)
	fmt.Printf("URL: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to Olist ERP API: %w", err)
	}

	// Add the Bearer token to the request header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Olist ERP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned non-200 status from Tiny ERP API: %s", resp.Status)
	}
	// Parse the response
	var apiResponse model.NFResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no NF found with number: %s", numero)
	}

	nfs := parseItensFromNFResponse(c, apiResponse, numero)
	return nfs, nil
}

func (c *NFController) GetByCPFCNPJ(ctx context.Context, cpfcnpj string) ([]*model.NF, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if nf, exists := c.nfs[cpfcnpj]; exists {
		return nf, nil
	}
	url := fmt.Sprintf("%s/notas?cpfCnpj=%s", OlistERPURL, cpfcnpj)
	fmt.Printf("URL: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to Olist ERP API: %w", err)
	}

	// Add the Bearer token to the request header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Olist ERP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned non-200 status from Tiny ERP API: %s", resp.Status)
	}
	// Parse the response
	var apiResponse model.NFResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}
	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no NF found with CPF/CNPJ: %s", cpfcnpj)
	}

	nfs := parseItensFromNFResponse(c, apiResponse, cpfcnpj)
	return nfs, nil
}

var nfSituacaoMap = map[string]string{
	"1":  "Pendente",
	"2":  "Emitida",
	"3":  "Cancelada",
	"4":  "Enviada Aguardando Recibo",
	"5":  "Rejeitada",
	"6":  "Autorizada",
	"7":  "Emitida Danfe",
	"8":  "Registrada",
	"9":  "Enviada Aguardando Protocolo",
	"10": "Denegada",
}

func parseItensFromNFResponse(c *NFController, apiResponse model.NFResponse, key string) []*model.NF {
	fmt.Println("Reached parseItensFromNFResponse")
	var nfs []*model.NF
	for _, item := range apiResponse.Itens {
		jsonData, _ := json.MarshalIndent(item, "", "  ")
		fmt.Printf("NF: %s\n", jsonData)

		situacaoDesc, ok := nfSituacaoMap[item.Situacao]
		if !ok {
			situacaoDesc = "Unknown"
		}

		id := item.ID
		url := fmt.Sprintf("%s/notas/%d/link", OlistERPURL, id)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			continue
		}
		defer resp.Body.Close()

		// Parse the response
		var secondaryResponse model.NFLink
		if err := json.NewDecoder(resp.Body).Decode(&secondaryResponse); err != nil {
			continue
		}
		nf := &model.NF{
			ID:          item.ID,
			Situacao:    situacaoDesc,
			Numero:      item.Numero,
			Serie:       item.Serie,
			ChaveAcesso: item.ChaveAcesso,
			DataEmissao: item.DataEmissao,
			Cliente: model.CustomerOrderEnpoint{
				CompanyOrderEndpoint: model.CompanyOrderEndpoint{
					TipoCadastro: "CLIENTE",
					ID:           item.Cliente.ID,
					Codigo:       model.Code(item.Cliente.Codigo),
					TipoPessoa:   item.Cliente.TipoPessoa,
					RazaoSocial:  item.Cliente.RazaoSocial,
					NomeFantasia: item.Cliente.NomeFantasia,
					CNPJCPF:      model.CNPJCPF(item.Cliente.CNPJCPF),
					Telefone:     model.Telefone(item.Cliente.Telefone),
					Email:        model.Email(item.Cliente.Email),
					Endereco: model.Endereco{
						Rua:         item.Cliente.Endereco.Rua,
						Numero:      item.Cliente.Endereco.Numero,
						Complemento: item.Cliente.Endereco.Complemento,
						Bairro:      item.Cliente.Endereco.Bairro,
						Municipio:   item.Cliente.Endereco.Municipio,
						Cep:         item.Cliente.Endereco.Cep,
						Uf:          item.Cliente.Endereco.Uf,
						Pais:        item.Cliente.Endereco.Pais,
					},
				},
			},
			Valor:              item.Valor,
			ValorProdutos:      item.ValorProdutos,
			CodigoRastreamento: item.CodigoRastreamento,
			UrlRastreamento:    item.UrlRastreamento,
			FretePorConta:      item.FretePorConta,
			Link:               secondaryResponse.Link,
		}
		nfs = append(nfs, nf)
	}
	c.nfs[key] = nfs
	return nfs
}
