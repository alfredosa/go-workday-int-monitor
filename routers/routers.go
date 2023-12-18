package routers

import (
	"net/http"

	"github.com/alfredosa/go-workday-int-monitor/handlers"
	"github.com/alfredosa/go-workday-int-monitor/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routers() *chi.Mux {
	api_config := handlers.WorkdayCredentials{
		ISU:      "ISU",
		Password: "password",
	}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.MiddlewareCors)

	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static")))
	r.Handle("/static/*", fileServer)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", api_config.HealthHandler)
	r.Mount("/api", apiRouter)

	appRouter := chi.NewRouter()
	appRouter.Route("/login", func(r chi.Router) {
		r.Get("/", api_config.LoginPage)
		r.Post("/", api_config.LoginPage)
	})
	appRouter.Get("/welcome", api_config.WelcomePage)
	r.Mount("/", appRouter)
	return r
}
