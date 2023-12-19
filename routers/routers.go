package routers

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/alfredosa/go-workday-int-monitor/handlers"
	"github.com/alfredosa/go-workday-int-monitor/middlewares"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func Routers() *chi.Mux {
	load_dotenv()
	setupDB()

	db, err := sql.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	api_config := handlers.Config{
		ISU:      "ISU",
		Password: "password",
		DB:       db,
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
	appRouter.Get("/loading", api_config.LoadingPage)
	appRouter.Get("/dashboard", api_config.DashboardPage)
	r.Mount("/", appRouter)
	return r
}

func setupDB() {
	if os.Getenv("LOCAL_DB") != "true" {
		return
	}
	if _, err := os.Stat("sqlite3.db"); err == nil {
		log.Println("sqlite3.db exists")
	} else if os.IsNotExist(err) {
		log.Println("sqlite3.db does not exist")
		log.Println("Creating sqlite3.db")
		file, err := os.Create("sqlite3.db")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	} else {
		log.Fatal(err)
	}
}

func load_dotenv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
