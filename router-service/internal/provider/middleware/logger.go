package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t1 := time.Now()
		reqID := middleware.GetReqID(ctx)

		preReqContent := slog.With(
			slog.String("requestTime", t1.Format(time.RFC3339)),
			slog.String("requestId", reqID),
			slog.String("method", r.Method),
			slog.String("endpoint", r.RequestURI),
			slog.String("protocol", r.Proto),
		)

		if r.RemoteAddr != "" {
			preReqContent.With(slog.String("ip", r.RequestURI))
		}

		preReqContent.Info("request started")

		defer func() {
			statusCode := 500
			if err := recover(); err != nil {
				slog.With(
					slog.String("requestId", reqID),
					slog.String("duration", time.Since(t1).String()),
					slog.Int("status", statusCode),
					slog.String("stacktrace", string(debug.Stack())),
				).Error("request finished with panic")
				panic(err)
			}
		}()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		status := ww.Status()

		log := slog.With(
			slog.String("requestId", reqID),
			slog.String("duration", time.Since(t1).String()),
			slog.Int("contentLength", ww.BytesWritten()),
			slog.Int("status", status),
		)

		statusCode := 400
		message := "request finished"

		if status >= statusCode {
			if err := ctx.Err(); err != nil {
				message += fmt.Sprintf(": %s", err.Error())
			}
			log.Error(message)
		} else {
			log.Info(message)
		}
	}

	return http.HandlerFunc(fn)
}
