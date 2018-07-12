package main

import (
	"net/http"

	"github.com/coderminer/microservice/showtimes/routes"
)

func main() {
	r := routes.NewRouter()
	http.ListenAndServe(":8002", r)
}
