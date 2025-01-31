package controller

import (
	"app/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type AccountReceivableController struct {
	accountsReceivable map[string][]*model.AccountReceivable
	customerController *CustomerController
	token              string
	mu                 sync.RWMutex
}

func NewAccountReceivableController(token string) *AccountReceivableController {
	return &AccountReceivableController{
		accountsReceivable: make(map[string][]*model.AccountReceivable),
		customerController: NewCustomerController(token),
		token:              token,
	}
}

func (c *AccountReceivableController) GetByName(ctx context.Context, name string) ([]*model.AccountReceivable, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if accountReceivable, exists := c.accountsReceivable[name]; exists {
		return accountReceivable, nil
	}

	customerName := c.customerController.GetNameFromPartialName(ctx, name)
	url := fmt.Sprintf("%s/contas-receber?nomeCliente=%s&situacao=aberto", OlistERPURL, customerName)
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
	var apiResponse model.AccountReceivableResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no account receivable found with name: %s", name)
	}

	// Cache the account receivable
	accounts := parseItensFromAccountReceivableResponse(c, apiResponse, name)

	return accounts, nil
}

func (c *AccountReceivableController) GetByNF(ctx context.Context, code string) ([]*model.AccountReceivable, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if accountReceivable, exists := c.accountsReceivable[code]; exists {
		return accountReceivable, nil
	}

	// url := fmt.Sprintf("%s/contas-receber?numeroDocumento=%s", OlistERPURL, code)
	// url := fmt.Sprintf("%s/contas-receber?numeroDocumento=%s&situacao=aberto", OlistERPURL, code)
	url := fmt.Sprintf("%s/contas-receber?situacao=aberto", OlistERPURL)
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
	var apiResponse model.AccountReceivableResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no account receivable found with code: %s", code)
	}

	accounts := parseItensFromAccountReceivableResponse(c, apiResponse, code)
	return accounts, nil
}

func (c *AccountReceivableController) GetByCPFCNPJ(ctx context.Context, cpfcnpj string) ([]*model.AccountReceivable, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if accountReceivable, exists := c.accountsReceivable[cpfcnpj]; exists {
		return accountReceivable, nil
	}

	// TODO: Need to get nomeCliente from CPF/CNPJ
	nomeCliente := c.customerController.GetNameFromCPFCNPJ(ctx, cpfcnpj)
	url := fmt.Sprintf("%s/contas-receber?nomeCliente=%s&situacao=aberto", OlistERPURL, nomeCliente)
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
	var apiResponse model.AccountReceivableResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no account receivable found with name: %s", cpfcnpj)
	}

	accounts := parseItensFromAccountReceivableResponse(c, apiResponse, cpfcnpj)
	return accounts, nil
}

func parseItensFromAccountReceivableResponse(c *AccountReceivableController, apiResponse model.AccountReceivableResponse, key string) []*model.AccountReceivable {
	fmt.Println("Reached parseItensFromAccountReceivableResponse")
	var accounts []*model.AccountReceivable
	for _, item := range apiResponse.Itens {
		if !strings.Contains(item.Doc, key) {
			continue
		}
		jsonData, _ := json.MarshalIndent(item, "", "  ")
		fmt.Printf("Account Receivable: %s\n", jsonData)
		accountReceivable := &model.AccountReceivable{
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
				}},
			Valor:      item.Valor,
			Doc:        item.Doc,
			Vencimento: item.Vencimento,
			Situacao:   item.Situacao,
		}

		accounts = append(accounts, accountReceivable)

	}
	c.accountsReceivable[key] = accounts
	return accounts
}

type AccountPayableController struct {
	accountsPayable map[string][]*model.AccountPayable
	token           string
	mu              sync.RWMutex
}

func NewAccountPayableController(token string) *AccountPayableController {
	return &AccountPayableController{
		accountsPayable: make(map[string][]*model.AccountPayable),
		token:           token,
	}
}

func (c *AccountPayableController) GetByName(ctx context.Context, name string) ([]*model.AccountPayable, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if accountPayable, exists := c.accountsPayable[name]; exists {
		return accountPayable, nil
	}

	url := fmt.Sprintf("%s/contas-pagar?nomeCliente=%s", OlistERPURL, name)
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
	var apiResponse model.AccountPayableResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no account payable found with name: %s", name)
	}

	// Cache the account payable
	accounts := parseItensFromAccountPayableResponse(c, apiResponse, name)

	return accounts, nil
}

func (c *AccountPayableController) GetByCPFCNPJ(ctx context.Context, cpfcnpj string) ([]*model.AccountPayable, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if accountPayable, exists := c.accountsPayable[cpfcnpj]; exists {
		return accountPayable, nil
	}
	url := fmt.Sprintf("%s/pedidos?cnpj=%s", OlistERPURL, cpfcnpj)
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
	var apiResponse model.AccountPayableResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no account payable found with name: %s", cpfcnpj)
	}

	accounts := parseItensFromAccountPayableResponse(c, apiResponse, cpfcnpj)
	return accounts, nil
}

func parseItensFromAccountPayableResponse(c *AccountPayableController, apiResponse model.AccountPayableResponse, key string) []*model.AccountPayable {
	fmt.Println("Reached parseItensFromAccountPayableResponse")
	var accounts []*model.AccountPayable
	for _, item := range apiResponse.Itens {
		jsonData, _ := json.MarshalIndent(item, "", "  ")
		fmt.Printf("Account Payable: %s\n", jsonData)
		/*
			order := &model.Order{
				ID:           item.ID,
				Situacao:     item.Situacao,
				NumeroPedido: item.NumeroPedido,
				Ecommerce:    item.Ecommerce,
				DataCriacao:  item.DataCriacao,
				DataPrevista: item.DataPrevista,
				Cliente: model.Customer{
					Company: model.Company{
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
					NomeVendedor: item.Vendedor.Nome,
				},
				Valor: item.Valor,
				Vendedor: model.Vendedor{
					ID:   item.Vendedor.ID,
					Nome: item.Vendedor.Nome,
				},
				Transportador: model.Transportador{
					ID:                 item.Transportador.ID,
					Nome:               item.Transportador.Nome,
					FretePorConta:      item.Transportador.FretePorConta,
					FormaEnvio:         model.FormaEnvio(item.Transportador.FormaEnvio),
					FormaFrete:         item.Transportador.FormaFrete,
					CodigoRastreamento: item.Transportador.CodigoRastreamento,
					UrlRastreamento:    item.Transportador.UrlRastreamento,
				},
			}
			orders = append(orders, order)
		*/
	}
	c.accountsPayable[key] = accounts
	return accounts
}
