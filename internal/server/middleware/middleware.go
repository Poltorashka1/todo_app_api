package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// MB Todo middleware structure

func PanicRecovery(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req)
			defer func() {
				if ok := recover(); ok != nil {
					log.Error(fmt.Sprintf("%v", ok))
				}
			}()
		})
	}
}

func HandlerExecutionTime(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, req)
			log.Info("Request: ", slog.String("method", req.Method), slog.String("path", req.RequestURI), slog.String("time", time.Since(start).String()))
		})
	}
}

// HandlerExecutionTimeV2 middleware for one method
// r.Get("/{id:\\d*}", middleware.HandlerExecutionTimeV2(http.HandlerFunc(server.Handlers.GetTaskHandler)))
func HandlerExecutionTimeV2(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("HandlerExecutionTimeV2")
		next.ServeHTTP(w, req)
		fmt.Println("HandlerExecutionTimeV2 end")
	}
}

// HandlerExecutionTimeV3 middleware for one method + logger
// r.Get("/{id:\\d*}", middleware.HandlerExecutionTimeV3(server.Log)(http.HandlerFunc(server.Handlers.GetTaskHandler)))
func HandlerExecutionTimeV3(log *slog.Logger) func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			log.Info(fmt.Sprintf("%s %s", req.Method, req.RequestURI))
			next.ServeHTTP(w, req)
		}
	}
}
