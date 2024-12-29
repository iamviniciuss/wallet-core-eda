package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/application/use_cases"
)

type WebAccountHandler struct {
	CreateAccountUseCase use_cases.CreateAccountUseCase
}

func NewWebAccountHandler(createAccountUseCase use_cases.CreateAccountUseCase) *WebAccountHandler {
	return &WebAccountHandler{
		CreateAccountUseCase: createAccountUseCase,
	}
}

func (h *WebAccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var dto use_cases.CreateAccountInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	output, err := h.CreateAccountUseCase.Execute(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
