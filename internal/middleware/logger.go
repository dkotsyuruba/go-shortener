package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/dkotsyuruba/go-shortener/internal/utils"
)

func LoggerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			startTime := time.Now()
			requestBody, _ := ReadRequestBody(req)
			recorder := &ResponseRecorder{
				ResponseWriter: w,
			}
			next.ServeHTTP(recorder, req)

			userID, _ := utils.GetUserID(req)
			logger.Info(
				"handled request",
				zap.String("uri", req.RequestURI),
				zap.String("method", req.Method),
				zap.Duration("duration", time.Since(startTime)),
				zap.Int("status_code", recorder.Status),
				zap.Int("size_bytes", recorder.Size),
				zap.String("user_id", userID),
				zap.ByteString("request_body", requestBody),
				zap.ByteString("response_body", recorder.Body.Bytes()),
			)
		})
	}
}

func ReadRequestBody(req *http.Request) ([]byte, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}

type ResponseRecorder struct {
	http.ResponseWriter
	Body      bytes.Buffer
	Status    int
	Size      int
	Wrote     bool
	HeaderMap http.Header
}

func (recorder *ResponseRecorder) Write(b []byte) (int, error) {
	recorder.Size += len(b)
	n, err := recorder.Body.Write(b)
	if err != nil {
		return n, err
	}
	return recorder.ResponseWriter.Write(b)
}
