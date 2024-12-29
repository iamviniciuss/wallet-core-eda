package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/application/use_cases"
)

type WebClientHandler struct {
	CreateClientUseCase use_cases.CreateClientUseCase
}

func NewWebClientHandler(createClientUseCase use_cases.CreateClientUseCase) *WebClientHandler {
	return &WebClientHandler{
		CreateClientUseCase: createClientUseCase,
	}
}

func (h *WebClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var dto use_cases.CreateClientInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	output, err := h.CreateClientUseCase.Execute(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
