package main

import (
	"log"
	"net/http"

	"kvstore/internal/handler"
	"kvstore/internal/router"
	"kvstore/internal/service"
	"kvstore/internal/store"
)

func main() {
	db := store.NewKVStore()
	svc := service.NewKVService(db)
	h := handler.NewKVHandler(svc)

	r := router.NewRouter(h)

	log.Println("Serving on port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
