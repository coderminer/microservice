package main

import (
	"net/http"

	"github.com/coderminer/microservice/movies/routes"
)

func main() {
	r := routes.NewRouter()
	http.ListenAndServe(":8001", r)
}
