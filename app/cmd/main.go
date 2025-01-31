package main

import (
	"app/internal/controller"
	httphandler "app/internal/handler/http"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nilsinho42/OlistERPMediator/auth/pkg/model"
)

var authToken string

func authenticate() error {
	authURL := fmt.Sprintf("%s/auth", os.Getenv("API_BASE_URL"))

	resp, err := http.Get(authURL)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth endpoint returned non-200 status: %s", resp.Status)
	}

	// Parse the response
	var token model.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	bytes, err := json.MarshalIndent(token, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token: %s", string(bytes))

	// Store the token in a global variable
	authToken = token.Key
	return nil
}

func main() {

	authenticate()
	// Initialize controllers
	supplierCtrl := controller.NewSupplierController(authToken)
	orderCtrl := controller.NewOrderController(authToken)
	customerCtrl := controller.NewCustomerController(authToken)
	accountReceivableCtrl := controller.NewAccountReceivableController(authToken)
	accountPayableCtrl := controller.NewAccountPayableController(authToken)

	// nfCtrl := controller.NewNFController()
	// financialCtrl := controller.NewFinancialController()
	// productCtrl := controller.NewProductController()

	// Initialize handler with controllers
	// h := httphandler.New(supplierCtrl, productCtrl, orderCtrl, customerCtrl, nfCtrl, financialCtrl)
	h := httphandler.New(supplierCtrl, customerCtrl, orderCtrl, accountReceivableCtrl, accountPayableCtrl)
	r := mux.NewRouter()

	r.HandleFunc("/v1/supplier", h.GetSupplier).Methods("GET")
	r.HandleFunc("/v1/order", h.GetOrder).Methods("GET")
	r.HandleFunc("/v1/customer", h.GetCustomer).Methods("GET")
	r.HandleFunc("/v1/boletos", h.GetAccountsReceivable).Methods("GET")

	// r.HandleFunc("/v1/product", h.GetProduct).Methods("GET")
	// r.HandleFunc("/v1/nf", h.GetNF).Methods("GET")
	// r.HandleFunc("/v1/financial", h.GetFinancial).Methods("GET")

	if err := http.ListenAndServe(":8082", r); err != nil {
		panic(err)
	}
}
