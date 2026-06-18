package transport

import (
	"context"
	"kvstore/internal/service"
	"kvstore/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedReplicationServer
	svc *service.KVService
}

func NewReplicationServer(svc *service.KVService) *Server {
	return &Server{svc: svc}
}

func (s *Server) GetSnapshot(req *pb.SnapshotRequest, stream pb.Replication_GetSnapshotServer) error {
	records, err := s.svc.GetSnapshotState()
	if err != nil {
		return err
	}
	for _, r := range records {
		mut := &pb.Mutation{
			Seq:   r.Seq,
			Op:    uint32(r.Op),
			Key:   r.Key,
			Value: r.Value,
		}
		if err := stream.Send(mut); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) ApplyMutation(ctx context.Context, req *pb.Mutation) (*pb.Ack, error) {
	err := s.svc.Set(req.Key, req.Value)
	if err != nil {
		return &pb.Ack{Ok: false}, err
	}
	return &pb.Ack{Ok: true}, nil
}

func StartGRPCServer(addr string, svc *service.KVService) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterReplicationServer(grpcServer, NewReplicationServer(svc))
	log.Printf("gRPC server listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
