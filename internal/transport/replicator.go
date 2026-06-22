package transport

// Replicator holds gRPC connections to all replicas and pushes every SET/DEL to them.

import (
	"context"
	"kvstore/internal/store"
	"kvstore/pb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Replicator struct {
	clients map[string]pb.ReplicationClient
	conns   []*grpc.ClientConn
}

func NewReplicator(addrs []string) *Replicator {
	clients := make(map[string]pb.ReplicationClient)
	var conns []*grpc.ClientConn

	for _, addr := range addrs {
		conn, err := grpc.NewClient(addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Printf("failed to connect to replica %s: %v", addr, err)
			continue
		}
		clients[addr] = pb.NewReplicationClient(conn)
		conns = append(conns, conn)
		log.Printf("Registered replica at %s (connects on first write)", addr)
	}

	return &Replicator{clients: clients, conns: conns}
}

// Takes mutation object,
// brodcast it to every client/connection stored on replicator
func (r *Replicator) Broadcast(mut store.Mutation) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.Mutation{
		Seq:   mut.Seq,
		Op:    uint32(mut.Op),
		Key:   mut.Key,
		Value: mut.Value,
	}

	for addr, client := range r.clients {
		ack, err := client.ApplyMutation(ctx, req)
		if err != nil {
			log.Printf("failed to replicate to %s: %v", addr, err)
			continue
		}
		log.Printf("replicated to %s: ok=%v", addr, ack.Ok)
	}
}

func (r *Replicator) Close() {
	for _, conn := range r.conns {
		conn.Close()
	}
}
