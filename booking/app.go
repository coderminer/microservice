package main

import (
	"net/http"

	"github.com/coderminer/microservice/booking/routes"
)

func main() {
	r := routes.NewRouter()
	http.ListenAndServe(":8003", r)
}
