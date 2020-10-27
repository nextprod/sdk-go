package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

// Handler represents extension handler function.
type Handler interface {
	Invoke(ctx context.Context, event []byte) ([]byte, error)
}

// extensionHandler represents generic extension handler function type
type extensionHandler func(context.Context, []byte) (interface{}, error)

// Invoke calls handler and serializes response.
func (h extensionHandler) Invoke(ctx context.Context, event []byte) ([]byte, error) {
	response, err := h(ctx, event)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func errorHandler(e error) extensionHandler {
	return func(ctx context.Context, event []byte) (interface{}, error) {
		return nil, e
	}
}

type serverOpt func(*server)

// NewHandler ...
func NewHandler(h interface{}) Handler {
	if h == nil {
		return errorHandler(fmt.Errorf("handler is nil"))
	}
	handler := reflect.ValueOf(h)
	htype := reflect.TypeOf(h)
	if htype.Kind() != reflect.Func {
		return errorHandler(fmt.Errorf("handler kind %s is not %s", htype.Kind(), reflect.Func))
	}
	return extensionHandler(func(ctx context.Context, payload []byte) (interface{}, error) {
		var args []reflect.Value
		args = append(args, reflect.ValueOf(ctx))
		eventType := htype.In(htype.NumIn() - 1)
		event := reflect.New(eventType)

		if err := json.Unmarshal(payload, event.Interface()); err != nil {
			return nil, err
		}
		args = append(args, event.Elem())

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
