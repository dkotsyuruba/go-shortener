package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			startTime := time.Now()
			resp := &ResponseProxy{w: w}
			next.ServeHTTP(resp, req)
			logger.Info(
				"handled request",
				zap.String("uri", req.RequestURI),
				zap.String("method", req.Method),
				zap.Duration("duration", time.Since(startTime)),
				zap.Int("status_code", resp.statusCode),
				zap.Int("size_bytes", resp.size),
			)
		})
	}
}

type ResponseProxy struct {
	statusCode int
	size       int
	w          http.ResponseWriter
}

func (rp *ResponseProxy) Write(data []byte) (int, error) {
	size, err := rp.w.Write(data)
	rp.size += size
	return size, err
}

func (rp *ResponseProxy) Header() http.Header {
	return rp.w.Header()
}

func (rp *ResponseProxy) WriteHeader(statusCode int) {
	rp.statusCode = statusCode
	rp.w.WriteHeader(statusCode)
}
