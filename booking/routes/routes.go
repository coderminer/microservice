package routes

import (
	"net/http"

	"github.com/coderminer/microservice/booking/controllers"
	"github.com/gorilla/mux"
)

type Route struct {
	Method      string
	Pattern     string
	Handler     http.HandlerFunc
	Middeleware mux.MiddlewareFunc
}

var routes []Route

func init() {
	register("POST", "/booking", controllers.CreateBooking, nil)
	register("GET", "/booking", controllers.GetAllBooking, nil)
	register("GET", "/booking/{name}", controllers.GetBookByName, nil)
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, Route{method, pattern, handler, middleware})
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range routes {
		r := router.Methods(route.Method).
			Path(route.Pattern)
		if route.Middeleware != nil {
			r.Handler(route.Middeleware(route.Handler))
		} else {
			r.Handler(route.Handler)
		}
	}
	return router
}
