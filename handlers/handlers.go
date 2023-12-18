package handlers

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

type WorkdayCredentials struct {
	ISU      string
	Password string
}

func (wdcreds *WorkdayCredentials) LoginPage(w http.ResponseWriter, r *http.Request) {
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

			http.Redirect(w, r, "/welcome", http.StatusFound)
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

func (wdcreds *WorkdayCredentials) WelcomePage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	if session.Value != "valid" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "Welcome, you have successfully logged in!")
}

func (wdcreds *WorkdayCredentials) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
