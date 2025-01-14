package routers

import (
	"fmt"
	"net/http"
	"forum/app/controllers"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Route struct {
	Handler http.HandlerFunc
	Method  string
}

type Router struct {
	Routes map[string]Route
}

func NewRouter() *Router {
	return &Router{Routes: make(map[string]Route)}
}

func (router *Router) AddRoute(path, method string, handler http.HandlerFunc, middlewares ...Middleware) {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	router.Routes[path] = Route{Handler: handler, Method: method}
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
	router.Routes[path] = Route{
		Handler: http.StripPrefix(path, fs).ServeHTTP,
		Method:  "GET",
	}
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
