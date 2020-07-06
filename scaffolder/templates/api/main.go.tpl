package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/lancer-kit/armory/api/render"
	"github.com/lancer-kit/armory/log"
	"github.com/lancer-kit/uwe/v2/presets/api"
	"github.com/sirupsen/logrus"

	"{{.project_name}}/config"
	"{{.project_name}}/info"
)

func GetServer(cfg *config.Cfg, logger *logrus.Entry) *api.Server {
	return api.NewServer(cfg.Api, getRouter(logger, cfg))
}

func getRouter(logger *logrus.Entry, cfg *config.Cfg) http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(log.NewRequestLogger(logger.Logger))

	if cfg.Api.EnableCORS {
		r.Use(getCORS().Handler)
	}

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	if cfg.Api.ApiRequestTimeout > 0 {
		t := time.Duration(cfg.Api.ApiRequestTimeout)
		r.Use(middleware.Timeout(t * time.Second))
	}

	r.Route("/dev", func(r chi.Router) {
		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			render.Success(w, info.App)
		})
		/*
			Here is needed to add http handlers logic for API worker
		 */

	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.ResultNotFound.Render(w)
	})

	return r
}
func getCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "jwt", "X-UID"},
		ExposedHeaders:   []string{"Link", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
