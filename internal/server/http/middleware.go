package internalhttp

import (
	"net/http"
	"time"

	"github.com/Dendyator/calendar/internal/logger" //nolint:depguard
)

type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	if !r.wroteHeader {
		r.status = statusCode
		r.ResponseWriter.WriteHeader(statusCode)
		r.wroteHeader = true
	}
}

func loggingMiddleware(logg *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(recorder, r)

			latency := time.Since(start)
			logg.Infof("%s %s [%s] \"%s %s %s\" %d %s %s",
				r.RemoteAddr,
				r.Header.Get("X-Forwarded-For"),
				start.Format("02/Jan/2006:15:04:05 -0700"),
				r.Method,
				r.URL.Path,
				r.Proto,
				recorder.status,
				latency,
				r.UserAgent(),
			)
		})
	}
}
