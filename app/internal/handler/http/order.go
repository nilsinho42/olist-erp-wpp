package http

import (
	"app/pkg/model"
	"net/http"
)

func (h *Handler) GetOrder(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := req.URL.Query()

	name := vars.Get("name")
	if name != "" {
		orders, err := h.orderController.GetByName(ctx, name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, orders, http.StatusOK)
		return
	}

	code := vars.Get("code")
	if code != "" {
		orders, err := h.orderController.GetByCode(ctx, code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, orders, http.StatusOK)
		return
	}

	cpfcnpj := vars.Get("cpfcnpj")
	if cpfcnpj != "" {
		cpfcnpj, err := model.ValidateCPFCNPJ(cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		orders, err := h.orderController.GetByCPFCNPJ(ctx, cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, orders, http.StatusOK)
		return
	}
}
