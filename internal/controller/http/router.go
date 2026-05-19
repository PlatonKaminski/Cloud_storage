package http

import (
	"log/slog"
	"net/http"
	"time"

	"cloud_storage/internal/controller/http/middleware"
	v1 "cloud_storage/internal/controller/http/v1"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type JWTValidator interface {
	Validate(token string) (string, error)
}

func NewRouter(handler *v1.Handler, jwt JWTValidator, log *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(loggerMiddleware(log))
	r.Use(chimiddleware.StripSlashes)

	r.Route("/auth", func(r chi.Router) {
		handler.RegisterAuth(r)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(jwt))
		handler.RegisterAPI(r)
	})

	r.Handle("/*", http.FileServer(http.Dir("./web")))

	return r
}

func loggerMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			log.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration", time.Since(start).String(),
				"request_id", chimiddleware.GetReqID(r.Context()),
			)
		})
	}
}
