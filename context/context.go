package context

import (
	"context"

	echo "github.com/labstack/echo/v4"
)

type echoCtx struct {
	context.Context
	ctx echo.Context
}

func (c *echoCtx) Value(key interface{}) interface{} {
	if k, ok := key.(string); ok {
		return c.ctx.Get(k)
	}

	return c.Context.Value(key)
}

func EchoContext(ctx echo.Context) context.Context {
	return &echoCtx{context.Background(), ctx}
}
