package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type ISU struct {
	ISU      string
	password string
}

func main() {
	http.HandleFunc("/", LoginPage)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/welcome", WelcomePage)

	// Serve static files from the "static" directory.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server on port 8080.
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
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

func WelcomePage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if session.Value != "valid" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "Welcome, you have successfully logged in!")
}
