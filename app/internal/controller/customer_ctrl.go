package controller

import (
	"app/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CustomerController struct {
	customers map[string][]*model.Customer
	token     string
	mu        sync.RWMutex
}

func NewCustomerController(token string) *CustomerController {
	return &CustomerController{
		customers: make(map[string][]*model.Customer),
		token:     token,
	}
}

func (c *CustomerController) GetByName(ctx context.Context, name string) ([]*model.Customer, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if customer, exists := c.customers[name]; exists {
		return customer, nil
	}
	url := fmt.Sprintf("%s/contatos?nome=%s", OlistERPURL, name)
	fmt.Printf("URL: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to Olist ERP API: %w", err)
	}

	// Add the Bearer token to the request header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Olist ERP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned non-200 status from Tiny ERP API: %s", resp.Status)
	}
	// Parse the response
	var apiResponse model.CustomerResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no customer found with name: %s", name)
	}

	// Cache the customer
	customers := parseItensFromCustomerResponse(c, apiResponse, name)

	return customers, nil
}

func (c *CustomerController) GetByCode(ctx context.Context, code string) ([]*model.Customer, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if customer, exists := c.customers[code]; exists {
		return customer, nil
	}
	url := fmt.Sprintf("%s/contatos?codigo=%s", OlistERPURL, code)
	fmt.Printf("URL: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to Olist ERP API: %w", err)
	}

	// Add the Bearer token to the request header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Olist ERP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned non-200 status from Tiny ERP API: %s", resp.Status)
	}
	// Parse the response
	var apiResponse model.CustomerResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no customer found with code: %s", code)
	}

	customers := parseItensFromCustomerResponse(c, apiResponse, code)
	return customers, nil
}

func (c *CustomerController) GetByCPFCNPJ(ctx context.Context, cpfcnpj string) ([]*model.Customer, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if customer, exists := c.customers[cpfcnpj]; exists {
		return customer, nil
	}
	url := fmt.Sprintf("%s/contatos?cpfCnpj=%s", OlistERPURL, cpfcnpj)
	fmt.Printf("URL: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to Olist ERP API: %w", err)
	}

	// Add the Bearer token to the request header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Olist ERP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned non-200 status from Tiny ERP API: %s", resp.Status)
	}
	// Parse the response
	var apiResponse model.CustomerResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no customer found with name: %s", cpfcnpj)
	}

	customers := parseItensFromCustomerResponse(c, apiResponse, cpfcnpj)
	return customers, nil
}

func (c *CustomerController) GetNameFromCPFCNPJ(ctx context.Context, cpfcnpj string) (customerName string) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	customers, err := c.GetByCPFCNPJ(ctx, cpfcnpj)
	if err != nil || len(customers) == 0 {
		return ""
	}
	return customers[0].CompanyCustomerEndpoint.RazaoSocial
}

func (c *CustomerController) GetNameFromPartialName(ctx context.Context, partialName string) (customerName string) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	customers, err := c.GetByName(ctx, partialName)
	if err != nil || len(customers) == 0 {
		return ""
	}
	return customers[0].CompanyCustomerEndpoint.RazaoSocial
}

func (c *CustomerController) GetCPFCNPJFromPartialName(ctx context.Context, partialName string) (customerName string) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	customers, err := c.GetByName(ctx, partialName)
	if err != nil || len(customers) == 0 {
		return ""
	}

	cpfcnpj, err := model.ValidateCPFCNPJ(string(customers[0].CompanyCustomerEndpoint.CNPJCPF))
	if err != nil {
		return ""
	}
	return cpfcnpj
}

func parseItensFromCustomerResponse(c *CustomerController, apiResponse model.CustomerResponse, key string) []*model.Customer {
	var customers []*model.Customer
	for _, item := range apiResponse.Itens {
		customer := &model.Customer{
			CompanyCustomerEndpoint: model.CompanyCustomerEndpoint{
				TipoCadastro: "CONSUMIDOR",
				ID:           item.ID,
				Codigo:       model.Code(item.Codigo),
				TipoPessoa:   item.TipoPessoa,
				RazaoSocial:  item.Nome,
				NomeFantasia: item.Fantasia,
				CNPJCPF:      model.CNPJCPF(item.CpfCnpj),
				Telefone:     model.Telefone(item.Telefone),
				Email:        model.Email(item.Email),
				Endereco: model.Endereco{
					Rua:         item.Endereco.Endereco,
					Numero:      item.Endereco.Numero,
					Complemento: item.Endereco.Complemento,
					Bairro:      item.Endereco.Bairro,
					Municipio:   item.Endereco.Municipio,
					Cep:         item.Endereco.Cep,
					Uf:          item.Endereco.Uf,
					Pais:        item.Endereco.Pais,
				},
			},
		}
		customers = append(customers, customer)
		c.customers[key] = customers
	}
	return customers
}
