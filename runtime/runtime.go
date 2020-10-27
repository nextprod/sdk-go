package runtime

import (
	"context"
	"encoding/json"
	"log"
	"os"

	pb "github.com/nextprod/sdk-go/pb"
)

func init() { log.SetOutput(os.Stdout) }

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
	var event interface{}
	if err := json.Unmarshal(in.GetEvent(), &event); err != nil {
		return &pb.InvokeReply{State: pb.State_Fail, Reason: err.Error()}, nil
	}
	_, err := s.invokeHandler.Invoke(ctx, event)
	if err != nil {
		return &pb.InvokeReply{State: pb.State_Fail, Reason: err.Error()}, nil
	}
	return &pb.InvokeReply{State: pb.State_Success}, nil
}

// withContext ...
func withContext(ctx context.Context) func(srv *server) {
	return func(srv *server) {
		srv.ctx = ctx
	}
}
