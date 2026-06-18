package router

import (
	"kvstore/internal/handler"
	"net/http"
)

func NewRouter(h *handler.KVHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /kv/{k}", func(w http.ResponseWriter, r *http.Request) { h.Get(w, r) })
	mux.HandleFunc("PUT /kv/{k}", func(w http.ResponseWriter, r *http.Request) { h.Set(w, r) })
	mux.HandleFunc("DELETE /kv/{k}", func(w http.ResponseWriter, r *http.Request) { h.Del(w, r) })

	return mux
}

func NewReplicaRouter(h *handler.KVHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /kv/{k}", func(w http.ResponseWriter, r *http.Request) { h.Get(w, r) })

	return mux
}
