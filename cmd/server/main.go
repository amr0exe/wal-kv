package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"kvstore/internal/handler"
	"kvstore/internal/router"
	"kvstore/internal/service"
	"kvstore/internal/store"
	"kvstore/internal/transport"
	typ "kvstore/internal/types"
)

func main() {
	node := flag.String("node", "primary", "node_role (primary|replica)")
	grpcPort := flag.String("port", ":5001", "gRPC replication port")
	httpPort := flag.String("http-port", ":8080", "HTTP server port")
	primaryAddr := flag.String("primary", "localhost:5001", "primary gRPC address (used by replica)")

	flag.Parse()

	role, err := ParseRole(*node)
	if err != nil {
		log.Fatal(err)
	}

	switch role {
	case typ.Primary:
		runPrimary(*httpPort, *grpcPort)
	case typ.Replica:
		runReplica(*httpPort, *primaryAddr)
	}
}

func runPrimary(httpAddr, grpcAddr string) {
	db, err := store.NewKVStore()
	if err != nil {
		log.Fatalf("failed to create store: %s", err.Error())
	}
	svc := service.NewKVService(db)

	go transport.StartGRPCServer(grpcAddr, svc)

	h := handler.NewKVHandler(svc)
	r := router.NewRouter(h)
	log.Printf("Primary HTTP server listening on %s", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, r))
}

func runReplica(httpAddr, primaryAddr string) {
	db := store.NewKVStoreInMemory()
	svc := service.NewKVService(db)

	client, err := transport.NewPrimaryClient(primaryAddr)
	if err != nil {
		log.Fatalf("failed to connect to primary: %s", err.Error())
	}
	defer client.Close()

	mutations, err := client.GetSnapshot()
	if err != nil {
		log.Fatalf("failed to get snapshot from primary: %s", err.Error())
	}
	db.LoadSnapshot(mutations)
	log.Printf("Replica loaded snapshot with %d entries", len(mutations))

	h := handler.NewKVHandler(svc)
	r := router.NewReplicaRouter(h)
	log.Printf("Replica HTTP server listening on %s", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, r))
}

func ParseRole(s string) (typ.Role, error) {
	switch s {
	case "primary":
		return typ.Primary, nil
	case "replica":
		return typ.Replica, nil
	default:
		return 0, fmt.Errorf("Invalid role: %s", s)
	}
}
