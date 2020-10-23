package runtime

import (
	"context"
	"log"

	pb "github.com/nextprod/sdk-go/pb"
)

const (
	port = ":50051"
)

// server is used to implement nex.RPCServer.
type server struct {
	ctx           context.Context
	invokeHandler Handler
	pb.UnimplementedRPCServer
}

// Invoke implements nex.RPCServer
func (s *server) Invoke(ctx context.Context, in *pb.InvokeRequest) (*pb.InvokeReply, error) {
	log.Printf("Received: %v", in)
	res, err := s.invokeHandler.Invoke(ctx, in.GetEvent())
	if err != nil {
		return nil, err
	}
	log.Printf("%v", res)
	return &pb.InvokeReply{}, nil
}

// withContext ...
func withContext(ctx context.Context) func(srv *server) {
	return func(srv *server) {
		srv.ctx = ctx
	}
}
