package contextutils

import "context"

type contextKey string

const RequestIDKey contextKey = "requestID"

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
