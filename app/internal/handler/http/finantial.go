package http

import (
	"net/http"
)

func (h *Handler) GetAccountsReceivable(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := req.URL.Query()

	name := vars.Get("name")
	if name != "" {
		accounts_receivable, err := h.accountReceivableController.GetByName(ctx, name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, accounts_receivable, http.StatusOK)
		return
	}

	cpfcnpj := vars.Get("cpfcnpj")
	if cpfcnpj != "" {
		cpfcnpj, err := validateCPFCNPJ(cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		accounts_receivable, err := h.accountReceivableController.GetByCPFCNPJ(ctx, cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, accounts_receivable, http.StatusOK)
		return
	}

	idNota := vars.Get("nf")
	if idNota != "" {
		accounts_receivable, err := h.accountReceivableController.GetByNF(ctx, idNota)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, accounts_receivable, http.StatusOK)
		return
	}

}

func (h *Handler) GetAccountsPayable(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := req.URL.Query()

	name := vars.Get("name")
	if name != "" {
		accounts_payable, err := h.accountPayableController.GetByName(ctx, name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, accounts_payable, http.StatusOK)
		return
	}

	cpfcnpj := vars.Get("cpfcnpj")
	if cpfcnpj != "" {
		cpfcnpj, err := validateCPFCNPJ(cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		accounts_payable, err := h.accountPayableController.GetByCPFCNPJ(ctx, cpfcnpj)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, accounts_payable, http.StatusOK)
		return
	}

}
