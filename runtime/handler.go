package runtime

import (
	"context"
	"encoding/json"
	"reflect"
)

// Handler represents extension handler function.
type Handler interface {
	Invoke(ctx context.Context, payload []byte) ([]byte, error)
}

// extensionHandler represents generic extension handler function type
type extensionHandler func(context.Context, []byte) (interface{}, error)

// Invoke calls handler and serializes response.
func (h extensionHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	response, err := h(ctx, payload)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type serverOpt func(*server)

// NewHandler ...
func NewHandler(h interface{}) Handler {
	handler := reflect.ValueOf(h)
	return extensionHandler(func(ctx context.Context, payload []byte) (interface{}, error) {
		var args []reflect.Value
		args = append(args, reflect.ValueOf(ctx))
		args = append(args, reflect.ValueOf(payload))
		res := handler.Call(args)

		var err error
		if len(res) > 0 {
			if errVal, ok := res[len(res)-1].Interface().(error); ok {
				err = errVal
			}
		}
		var val interface{}
		if len(res) > 1 {
			val = res[0].Interface()
		}
		return val, err
	})
}
