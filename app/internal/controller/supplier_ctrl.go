package controller

import (
	"app/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type SupplierController struct {
	suppliers map[string][]*model.Supplier
	token     string
	mu        sync.RWMutex
}

func NewSupplierController(token string) *SupplierController {
	return &SupplierController{
		suppliers: make(map[string][]*model.Supplier),
		token:     token,
	}
}

func (c *SupplierController) GetByName(ctx context.Context, name string) ([]*model.Supplier, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if supplier, exists := c.suppliers[name]; exists {
		return supplier, nil
	}
	url := fmt.Sprintf("%s/contatos?nome=%s", OlistERPURL, name)
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
	var apiResponse model.SupplierResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no supplier found with name: %s", name)
	}

	// Cache the supplier
	suppliers := parseItensFromSupplierResponse(c, apiResponse, name)

	return suppliers, nil
}

func (c *SupplierController) GetByCode(ctx context.Context, code string) ([]*model.Supplier, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if supplier, exists := c.suppliers[code]; exists {
		return supplier, nil
	}
	url := fmt.Sprintf("%s/contatos?codigo=%s", OlistERPURL, code)
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
	var apiResponse model.SupplierResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no supplier found with code: %s", code)
	}

	suppliers := parseItensFromSupplierResponse(c, apiResponse, code)
	return suppliers, nil
}

func (c *SupplierController) GetByCPFCNPJ(ctx context.Context, cpfcnpj string) ([]*model.Supplier, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if supplier, exists := c.suppliers[cpfcnpj]; exists {
		return supplier, nil
	}
	url := fmt.Sprintf("%s/contatos?cpfCnpj=%s", OlistERPURL, cpfcnpj)
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
	var apiResponse model.SupplierResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no supplier found with name: %s", cpfcnpj)
	}

	suppliers := parseItensFromSupplierResponse(c, apiResponse, cpfcnpj)
	return suppliers, nil
}

func parseItensFromSupplierResponse(c *SupplierController, apiResponse model.SupplierResponse, key string) []*model.Supplier {
	var suppliers []*model.Supplier
	for _, item := range apiResponse.Itens {
		supplier := &model.Supplier{
			Company: model.Company{
				TipoCadastro: "FORNECEDOR",
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
		suppliers = append(suppliers, supplier)
		c.suppliers[key] = suppliers
	}
	return suppliers
}
