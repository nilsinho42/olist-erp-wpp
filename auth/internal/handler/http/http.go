package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"auth/internal/controller"
	"auth/pkg/model"
)

type Handler struct {
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetToken(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	token, err := h.ctrl.Get(ctx)
	if err != nil {
		log.Printf("Error getting token: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(token); err != nil {
		log.Printf("Error encoding token: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PutToken(w http.ResponseWriter, req *http.Request) {
	// example curl to validate this:
	// curl -X PUT "http://localhost:8081/auth?key=34adade1-6ac4-4a5a-a394-2c47177a9311.95c5eb2f-e8a8-4f48-8bf2-fa2882f6c607.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb"
	vars := req.URL.Query()

	ctx := req.Context()
	ctx = context.WithValue(ctx, model.ContextKey, vars.Get("key"))
	if err := h.ctrl.Put(ctx); err != nil {
		log.Printf("Error putting token: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
