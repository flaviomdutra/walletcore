package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/flaviomdutra/walletcore/internal/usecase/create_client"
)

type WebClientHandler struct {
	CreateClientUseCase create_client.CreateClientUseCase
}

func NewWebClientHandler(createClientUseCase create_client.CreateClientUseCase) *WebClientHandler {
	return &WebClientHandler{
		CreateClientUseCase: createClientUseCase,
	}
}

func (h *WebClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var dto create_client.CreateClientInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "invalid JSON request body", http.StatusBadRequest)
		return
	}

	output, err := h.CreateClientUseCase.Execute(dto)
	if err != nil {
		log.Printf("create client failed: %v", err)
		http.Error(w, "could not create client", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		log.Printf("encode create client response failed: %v", err)
		return
	}
}
