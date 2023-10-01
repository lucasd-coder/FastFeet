package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/lucasd-coder/router-service/pkg/logger"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t1 := time.Now()
		reqID := middleware.GetReqID(ctx)

		preReqContent := map[string]interface{}{
			"requestTime": t1.Format(time.RFC3339),
			"requestId":   reqID,
			"method":      r.Method,
			"endpoint":    r.RequestURI,
			"protocol":    r.Proto,
		}

		if r.RemoteAddr != "" {
			preReqContent["ip"] = r.RemoteAddr
		}

		log := logger.FromContext(ctx).WithFields(preReqContent)
		log.Info("request started")

		defer func() {
			statusCode := 500
			if err := recover(); err != nil {
				log.WithFields(logger.Fields{
					"requestId":  reqID,
					"duration":   time.Since(t1).String(),
					"status":     statusCode,
					"stacktrace": string(debug.Stack()),
				}).Error("request finished with panic")
				panic(err)
			}
		}()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		status := ww.Status()

		postReqContent := map[string]interface{}{
			"requestId":     reqID,
			"duration":      time.Since(t1).String(),
			"contentLength": ww.BytesWritten(),
			"status":        status,
		}
		log = logger.FromContext(ctx).WithFields(postReqContent)

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
