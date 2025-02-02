package http

import (
	"net/http"
)

func (h *Handler) GetProduct(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := req.URL.Query()

	name := vars.Get("name")
	if name != "" {
		products, err := h.productController.GetByName(ctx, name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		WriteResponse(w, products, http.StatusOK)
		return
	}
}
