package http

import (
	"app/internal/controller"
	"encoding/json"
	"net/http"
)

type Handler struct {
	supplierController          *controller.SupplierController
	customerController          *controller.CustomerController
	orderController             *controller.OrderController
	accountReceivableController *controller.AccountReceivableController
	accountPayableController    *controller.AccountPayableController
	nfController                *controller.NFController
	productController           *controller.ProductController
	// financialController *controller.FinancialController
}

func New(
	supplierCtrl *controller.SupplierController,
	customerCtrl *controller.CustomerController,
	orderCtrl *controller.OrderController,
	accountReceivableCtrl *controller.AccountReceivableController,
	accountPayableCtrl *controller.AccountPayableController,
	nfCtrl *controller.NFController,
	productCtrl *controller.ProductController,
	// financialCtrl *controller.FinancialController,
) *Handler {
	return &Handler{
		supplierController:          supplierCtrl,
		customerController:          customerCtrl,
		orderController:             orderCtrl,
		accountReceivableController: accountReceivableCtrl,
		accountPayableController:    accountPayableCtrl,
		nfController:                nfCtrl,
		productController:           productCtrl,
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
