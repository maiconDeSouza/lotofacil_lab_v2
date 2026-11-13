package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/service"
)

type HandlerInterface interface {
	Home(w http.ResponseWriter, r *http.Request)
	GetContest(w http.ResponseWriter, r *http.Request)
	GetNumberContest(w http.ResponseWriter, r *http.Request)
	GenerateDailyContest(w http.ResponseWriter, r *http.Request)
	GenerateDailyGame(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	service service.ServiceInterface
}

func NewHandler(service service.ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func (h *Handler) GetContest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	contest, errappRes := h.service.GetContest()

	if errappRes != nil {
		w.WriteHeader(errappRes.Code)
		json.NewEncoder(w).Encode(errappRes)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contest)

}

func (h *Handler) GetNumberContest(w http.ResponseWriter, r *http.Request) {
	c := r.PathValue("concurso")
	w.Header().Set("Content-Type", "application/json")
	contest, errappRes := h.service.GetNumberContest(c)

	if errappRes != nil {
		w.WriteHeader(errappRes.Code)
		json.NewEncoder(w).Encode(errappRes)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contest)
}

func (h *Handler) GenerateDailyContest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	contest, errRes := h.service.GenerateDailyContest()

	if errRes != nil {
		w.WriteHeader(errRes.Code)
		json.NewEncoder(w).Encode(errRes)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contest)
}

func (h *Handler) GenerateDailyGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	games, errRes := h.service.GenerateDailyGames()
	if errRes != nil {
		w.WriteHeader(errRes.Code)
		json.NewEncoder(w).Encode(errRes)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(games)
}
