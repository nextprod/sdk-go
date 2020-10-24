package runtime

import (
	"context"
	"fmt"
	"testing"

	"github.com/nextprod/sdk-go/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func TestRuntime(t *testing.T) {
	ctx := context.Background()
	srv := &server{
		invokeHandler: NewHandler(func(ctx context.Context, payload []byte) ([]byte, error) {
			fmt.Println("asd")
			return nil, nil
		}),
	}
	_, err := srv.Invoke(ctx, &pb.InvokeRequest{
		Event: []byte{'1'},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestIngr(t *testing.T) {
	go func() {
		Start(func(ctx context.Context, payload []byte) ([]byte, error) {
			return nil, nil
		})
	}()
	// Invoke extension
	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewRPCClient(conn)
	reply, err := c.Invoke(context.Background(), &pb.InvokeRequest{Event: []byte(string("hello from agent"))})
	if err != nil {
		t.Fatal(err)
	}
	logrus.Infof("Relply from extension: %s", reply)
}
