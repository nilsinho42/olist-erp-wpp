package http

import (
	"app/pkg/model"
	"net/http"
)

func (h *Handler) GetCustomer(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := req.URL.Query()

	name := vars.Get("name")
	if name != "" {
		customers, err := h.customerController.GetByName(ctx, name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			WriteResponse(w, customers, http.StatusOK)
		}
		return
	}

	code := vars.Get("code")
	if code != "" {
		customers, err := h.customerController.GetByCode(ctx, code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			WriteResponse(w, customers, http.StatusOK)
		}
		return
	}

	cpfcnpj := vars.Get("cpfcnpj")
	if cpfcnpj != "" {
		cpfcnpj, err := model.ValidateCPFCNPJ(cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		customers, err := h.customerController.GetByCPFCNPJ(ctx, cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			WriteResponse(w, customers, http.StatusOK)
		}
		return
	}
}
