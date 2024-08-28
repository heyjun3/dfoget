package server

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
)

const (
	AuthCookieName = "dforget"
)

func NewAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			cookie := req.Header().Get("Cookie")
			name, val, _ := strings.Cut(cookie, "=")
			if name != AuthCookieName {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					fmt.Errorf("expect cookie name is %s", AuthCookieName),
				)
			}
			slog.InfoContext(ctx, "cookie", "value", val)
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
