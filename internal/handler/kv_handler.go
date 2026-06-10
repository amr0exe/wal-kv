package handler

import (
	"encoding/json"
	"kvstore/internal/service"
	"net/http"
)

// Handler layer acts between the service and requests
// Generally handles request/response related operations

type KVHandler struct {
	svc *service.KVService
}

func NewKVHandler(svc *service.KVService) *KVHandler {
	return &KVHandler{
		svc: svc,
	}
}

type SetReq struct {
	Value string `json:"value"`
}

func (h *KVHandler) Set(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	key := r.PathValue("k")

	var req SetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := h.svc.Set(key, req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("ok"))
}

func (h *KVHandler) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	key := r.PathValue("k")

	err, value := h.svc.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(value))
}

func (h *KVHandler) Del(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	key := r.PathValue("k")

	err := h.svc.Del(key)
	if err != nil {
		http.Error(w, "failed delete operation", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNoContent)
}
