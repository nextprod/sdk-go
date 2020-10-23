package runtime

import (
	"context"
	"testing"

	"github.com/nextprod/sdk-go/pb"
)

func TestRuntime(t *testing.T) {
	ctx := context.Background()
	srv := &server{
		invokeHandler: NewHandler(func(ctx context.Context, payload []byte) ([]byte, error) {
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
