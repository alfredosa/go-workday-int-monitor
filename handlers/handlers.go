package handlers

import (
	"database/sql"
	"net/http"
	"text/template"
	"time"
)

type Config struct {
	ISU      string
	Password string
	DB       *sql.DB
}

func (config *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	var errorMessage string

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username != "admin" || password != "password" {
			errorMessage = "Invalid credentials"

			http.SetCookie(w, &http.Cookie{
				Name:    "session",
				Value:   "invalid",
				Expires: time.Now().Add(1 * time.Hour),
			})
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:    "session",
				Value:   "valid",
				Expires: time.Now().Add(1 * time.Hour),
			})

			http.Redirect(w, r, "/loading", http.StatusFound)
			return
		}
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]string{
		"ErrorMessage": errorMessage,
	})
}

func (config *Config) DashboardPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	checkCookies(w, r, "/login")
	tmpl.Execute(w, map[string]string{
		"ErrorMessage": "",
	})
}

func (config *Config) LoadingPage(w http.ResponseWriter, r *http.Request) {
	checkCookies(w, r, "/login")

	tmpl, err := template.ParseFiles("templates/loading_page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]string{
		"ErrorMessage": "",
	})
}

func (config *Config) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func checkCookies(w http.ResponseWriter, r *http.Request, redirect string) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}
	if cookie.Value != "valid" {
		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}
}
