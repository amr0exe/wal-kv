package transport

// Client used by replica to fetch a full-state snapshot from the primary on boot.

import (
	"context"
	"io"
	"kvstore/internal/store"
	ty "kvstore/internal/types"
	"kvstore/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PrimaryClient struct {
	client pb.ReplicationClient
	conn   *grpc.ClientConn
}

func NewPrimaryClient(addr string) (*PrimaryClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &PrimaryClient{
		client: pb.NewReplicationClient(conn),
		conn:   conn,
	}, nil
}

func (c *PrimaryClient) Close() error {
	return c.conn.Close()
}

// GetSnapshot[Replica] method, sends empty SnapshotRequest,
// Asking for stream of mutations
func (c *PrimaryClient) GetSnapshot() ([]store.Mutation, error) {
	stream, err := c.client.GetSnapshot(context.Background(), &pb.SnapshotRequest{})
	if err != nil {
		return nil, err
	}

	var mutations []store.Mutation
	for {
		mut, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		mutations = append(mutations, store.Mutation{
			Seq:   mut.Seq,
			Op:    ty.OpType(mut.Op),
			Key:   mut.Key,
			Value: mut.Value,
		})
	}
	return mutations, nil
}
