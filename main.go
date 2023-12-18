package main

import (
	"log"
	"net/http"

	"github.com/alfredosa/go-workday-int-monitor/routers"
)

func main() {
	const port string = ":8080"
	r := routers.Routers()

	log.Printf("Serving on Port: %s\n", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}
