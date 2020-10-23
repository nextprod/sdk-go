package runtime

import (
	"context"
	"fmt"
	"net"

	pb "github.com/nextprod/sdk-go/pb"
	"google.golang.org/grpc"
)

func startGRPC(ctx context.Context, handler Handler) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterRPCServer(s, &server{ctx: ctx, invokeHandler: handler})
	if err := s.Serve(ln); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}
	return nil
}

func init() { rpcStartFunc.f = startGRPC }
