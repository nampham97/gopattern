package logger

import (
	"context"

	"GoPattern/middleware/contextutils"

	"go.uber.org/zap"
)

type ContextLogger struct {
	ctx context.Context
	*zap.Logger
}

func FromContext(ctx context.Context) *ContextLogger {
	reqID := contextutils.GetRequestID(ctx)
	return &ContextLogger{
		ctx: ctx,
		Logger: Log.With(
			zap.String("request_id", reqID),
		),
	}
}
