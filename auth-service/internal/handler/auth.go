package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/esmrzv/microservice-shop/auth-service/internal/service"
	"github.com/esmrzv/microservice-shop/auth-service/internal/service/dto"
)


type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(h service.AuthService) AuthHandler {
	return &authHandler{service: h}
}


func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req dto.RegisterRequest

	err :=json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	response, err := h.service.Register(r.Context(), req)
	if err != nil{
		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request){
	var req dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	response, err := h.service.Login(r.Context(), req)
	if err != nil{
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "failed to login user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}