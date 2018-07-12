package main

import (
	"net/http"

	"github.com/coderminer/microservice/users/routes"
)

func main() {
	r := routes.NewRouter()
	http.ListenAndServe(":8000", r)
}
