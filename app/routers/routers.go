package routers

import (
	"fmt"
	"net/http"

	"forum/app/controllers"
)

// Middleware defines a function type that wraps an http.HandlerFunc.
// It allows adding pre- or post-processing logic around the original handler.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Route represents an individual route in the application.
type Route struct {
	Method  string           // The HTTP method required for the route (e.g., "GET", "POST").
	Handler http.HandlerFunc // The function to execute when the route is matched.
}

// Router manages all routes for the application.
type Router struct {
	Routes map[string]*Route // A map of paths to their corresponding Route.
}

// Intantiate a new route:
func newRoute(method string, handler http.HandlerFunc) *Route {
	route := new(Route)
	route.Method = method
	route.Handler = handler
	return route
}

// Instantiate a new router:
func NewRouter() *Router {
	router := new(Router)
	// mapping all routes where the key is the path and the value is a the route:
	// key: /path  value: route(method->handler);
	router.Routes = make(map[string]*Route)
	return router
}

// AddRoute adds a new route to the router with optional middleware.
// Parameters:
// - path: The URL path to match (e.g., "/login").
// - method: The HTTP method required (e.g., "GET", "POST").
// - handler: The function to execute for the route.
// - middlewares: A variadic argument allowing multiple middleware to be applied.
func (router *Router) AddRoute(path, method string, handler http.HandlerFunc, middlewares ...Middleware) {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	router.Routes[path] = newRoute(method, handler)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route, exists := router.Routes[req.URL.Path]
	if exists && route.Method == req.Method {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		route.Handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func (router *Router) AddStaticRoute(path, dir string) {
	fs := http.FileServer(http.Dir(dir))
	router.Routes[path] =  newRoute("GET" ,http.StripPrefix(path, fs).ServeHTTP)
}

func (router *Router) ListRoutes() {
	for path, route := range router.Routes {
		fmt.Printf("Path: %s, Method: %s\n", path, route.Method)
	}
}

func (router *Router) RouteHandler() {
	router.AddRoute("/", "GET", controllers.RegisterUserHandler)
	router.AddRoute("/login", "POST", controllers.LoginHandler)
	router.AddRoute("/posts", "GET", controllers.PostsHandler)
	// Add other routes as needed
}
