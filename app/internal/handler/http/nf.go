package http

import (
	"app/pkg/model"
	"net/http"
)

func (h *Handler) GetNF(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := req.URL.Query()

	name := vars.Get("name")
	if name != "" {
		nfs, err := h.nfController.GetByName(ctx, name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, nfs, http.StatusOK)
		return
	}

	numero := vars.Get("numero")
	if numero != "" {
		nfs, err := h.nfController.GetByNumero(ctx, numero)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, nfs, http.StatusOK)
		return
	}

	cpfcnpj := vars.Get("cpfcnpj")
	if cpfcnpj != "" {
		cpfcnpj, err := model.ValidateCPFCNPJ(cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		nfs, err := h.nfController.GetByCPFCNPJ(ctx, cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, nfs, http.StatusOK)
		return
	}
}
