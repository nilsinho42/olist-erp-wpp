package http

import (
	"app/internal/controller"
	"fmt"
	"net/http"
)

type Handler struct {
	supplierController *controller.SupplierController
	// productController   *controller.ProductController
	// orderController     *controller.OrderController
	// customerController  *controller.CustomerController
	// nfController        *controller.NFController
	// financialController *controller.FinancialController
}

func New(
	supplierCtrl *controller.SupplierController,
	// productCtrl *controller.ProductController,
	// orderCtrl *controller.OrderController,
	// customerCtrl *controller.CustomerController,
	// nfCtrl *controller.NFController,
	// financialCtrl *controller.FinancialController,
) *Handler {
	return &Handler{
		supplierController: supplierCtrl,
		// productController:   productCtrl,
		// orderController:     orderCtrl,
		// customerController:  customerCtrl,
		// nfController:        nfCtrl,
		// financialController: financialCtrl,
	}
}

func (h *Handler) GetSupplier(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := req.URL.Query()

	name := vars.Get("name")
	fmt.Println(name)
	if name != "" {
		supplier, err := h.supplierController.GetByName(ctx, name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println(supplier)
		fmt.Println(err)
		return
	}

	code := vars.Get("code")
	if code == "" {

		return
	}

	vendor := vars.Get("vendor")
	if vendor == "" {

		return
	}

	cpfcnpj := vars.Get("cpfcnpj")
	if cpfcnpj == "" {

		return
	}

}
