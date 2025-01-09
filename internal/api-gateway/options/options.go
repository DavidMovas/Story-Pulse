package options

import (
	"context"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"net/http"
)

func CustomErrorHandler(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	grpcErr := status.Convert(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(runtime.HTTPStatusFromCode(grpcErr.Code()))
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": grpcErr.Message(),
	})
}

type CustomMarshaler struct {
	*runtime.JSONPb
}

func (c CustomMarshaler) Marshal(v interface{}) ([]byte, error) {
	if entityMap, ok := v.(map[string]interface{}); ok && len(entityMap) == 1 {
		for _, value := range entityMap {
			return json.Marshal(value)
		}
	}

	return c.JSONPb.Marshal(v)
}
