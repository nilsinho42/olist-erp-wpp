package http

import (
	"app/pkg/model"
	"net/http"
)

func (h *Handler) GetSupplier(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := req.URL.Query()

	name := vars.Get("name")
	if name != "" {
		suppliers, err := h.supplierController.GetByName(ctx, name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, suppliers, http.StatusOK)
		return
	}

	code := vars.Get("code")
	if code != "" {
		suppliers, err := h.supplierController.GetByCode(ctx, code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, suppliers, http.StatusOK)
		return
	}

	cpfcnpj := vars.Get("cpfcnpj")
	if cpfcnpj != "" {
		cpfcnpj, err := model.ValidateCPFCNPJ(cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		suppliers, err := h.supplierController.GetByCPFCNPJ(ctx, cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, suppliers, http.StatusOK)
		return
	}
}
