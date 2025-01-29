package http

import (
	"app/internal/controller"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Handler struct {
	supplierController *controller.SupplierController
	customerController *controller.CustomerController
	orderController    *controller.OrderController

	// productController   *controller.ProductController
	// nfController        *controller.NFController
	// financialController *controller.FinancialController
}

func New(
	supplierCtrl *controller.SupplierController,
	customerCtrl *controller.CustomerController,
	orderCtrl *controller.OrderController,

	// productCtrl *controller.ProductController,
	// nfCtrl *controller.NFController,
	// financialCtrl *controller.FinancialController,
) *Handler {
	return &Handler{
		supplierController: supplierCtrl,
		customerController: customerCtrl,
		orderController:    orderCtrl,
		// productController:   productCtrl,
		// nfController:        nfCtrl,
		// financialController: financialCtrl,
	}
}

func WriteResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
	// fmt.Printf("Suppliers: %s\n", jsonData)
}

func validateCPFCNPJ(input string) (string, error) {
	// If CNPJ: 15.049.188/0001-30 => 15.049.188/0001-30, nil
	// If CNPJ: 15049188000130 => 15.049.188/0001-30, nil
	// If CPF: 123.456.789-09 => 123.456.789-09, nil
	// If CPF: 12345678909 => 123.456.789-09, nil
	// any other format:  given, error

	// If no ., / or - AND length is 14, split in 3 parts and add . and / or - in the right places
	// If no . or - AND length is 11, split in 4 parts and add . and - in the right places

	input = strings.ReplaceAll(input, ".", "")
	input = strings.ReplaceAll(input, "-", "")
	input = strings.ReplaceAll(input, "/", "")
	if len(input) == 14 {
		fmt.Printf("%s.%s.%s/%s-%s\n", input[:2], input[2:5], input[5:8], input[8:12], input[12:14])
		return fmt.Sprintf("%s.%s.%s/%s-%s", input[:2], input[2:5], input[5:8], input[8:12], input[12:14]), nil
	} else if len(input) == 11 {
		return fmt.Sprintf("%s.%s.%s-%s", input[:3], input[3:6], input[6:9], input[9:11]), nil
	}
	return input, fmt.Errorf("invalid CPF/CNPJ format")

}
