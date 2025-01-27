package controller

import (
	"app/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

var OlistERPURL = os.Getenv("API_OLIST_URL")

type SupplierController struct {
	suppliers map[string]*model.Supplier
	token     string
	mu        sync.RWMutex
}

func NewSupplierController(token string) *SupplierController {
	return &SupplierController{
		suppliers: make(map[string]*model.Supplier),
		token:     token,
	}
}

func (c *SupplierController) GetByName(ctx context.Context, name string) (*model.Supplier, error) {
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
	req.Header.Add("Authorization: ", fmt.Sprintf("Bearer %s", c.token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Olist ERP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Tiny ERP API returned non-200 status: %s", resp.Status)
	}

	// Parse the response
	var apiResponse struct {
		Retorno struct {
			Contato model.Supplier `json:"contato"`
		} `json:"retorno"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Tiny ERP API response: %w", err)
	}

	// Cache the supplier
	c.suppliers[name] = &apiResponse.Retorno.Contato

	return &apiResponse.Retorno.Contato, nil
}
