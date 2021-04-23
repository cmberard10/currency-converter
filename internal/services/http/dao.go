package http

import (
	"context"
	"net/http"
)

//DAO ...
type DAO interface {
	Do(ctx context.Context, request *http.Request) ([]byte, int, error)
	WithTimeout(ctx context.Context) (context.Context, context.CancelFunc)
}
