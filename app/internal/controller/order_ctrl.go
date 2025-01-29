package controller

import (
	"app/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type OrderController struct {
	orders map[string][]*model.Order
	token  string
	mu     sync.RWMutex
}

func NewOrderController(token string) *OrderController {
	return &OrderController{
		orders: make(map[string][]*model.Order),
		token:  token,
	}
}

func (c *OrderController) GetByName(ctx context.Context, name string) ([]*model.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if order, exists := c.orders[name]; exists {
		return order, nil
	}
	url := fmt.Sprintf("%s/pedidos?nomeCliente=%s", OlistERPURL, name)
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
	/* TODO REMOVE
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Tiny ERP API response body: %w", err)
	}
	fmt.Printf("Response: %s\n", string(bodyBytes))
	*/
	// Parse the response
	var apiResponse model.OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no order found with name: %s", name)
	}

	// Cache the order
	orders := parseItensFromOrderResponse(c, apiResponse, name)

	return orders, nil
}

func (c *OrderController) GetByCode(ctx context.Context, code string) ([]*model.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if order, exists := c.orders[code]; exists {
		return order, nil
	}
	url := fmt.Sprintf("%s/pedidos?codigoCliente=%s", OlistERPURL, code)
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
	var apiResponse model.OrderResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no order found with code: %s", code)
	}

	orders := parseItensFromOrderResponse(c, apiResponse, code)
	return orders, nil
}

func (c *OrderController) GetByCPFCNPJ(ctx context.Context, cpfcnpj string) ([]*model.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if order, exists := c.orders[cpfcnpj]; exists {
		return order, nil
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
	var apiResponse model.OrderResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	if len(apiResponse.Itens) == 0 {
		return nil, fmt.Errorf("no order found with name: %s", cpfcnpj)
	}

	orders := parseItensFromOrderResponse(c, apiResponse, cpfcnpj)
	return orders, nil
}

func parseItensFromOrderResponse(c *OrderController, apiResponse model.OrderResponse, key string) []*model.Order {
	fmt.Println("Reached parseItensFromOrderResponse")
	var orders []*model.Order
	for _, item := range apiResponse.Itens {
		jsonData, _ := json.MarshalIndent(item, "", "  ")
		fmt.Printf("Order: %s\n", jsonData)

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
	}
	c.orders[key] = orders
	return orders
}
