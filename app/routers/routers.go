package routers

import (
	"net/http"
	"forum/app/controllers"
)

// Declare a structure to represent the router entity:
type Router struct {
	Routes map[string]http.HandlerFunc
}

// Instantiate a new router:
func NewRouter() *Router {
	router := new(Router)
	router.Routes = make(map[string]http.HandlerFunc)
	return router
}

// Add a route and its handler to the router:
func (router *Router) AddRoute(path string, handler http.HandlerFunc) {
	router.Routes[path] = handler
}

// ServeHTTP is called when a request is made to the server
func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Match the URL path to a handler
	handler, exists := router.Routes[req.URL.Path]
	if exists {
		handler(w, req) // Call the handler if the route exists
	} else {
		http.NotFound(w, req) // Return a 404 if no match is found
	}
}


// Create a new handler:
func (router *Router) RouteHandler() {
	router.AddRoute("/", controllers.HomeHandler)
	router.AddRoute("/user", controllers.GetAllUsersHandler)
}