package routers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"forum/app/controllers"
)

// HandlerFunc defines the type for handler functions
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Router stores routes and their corresponding handlers
type Router struct {
	routes      map[string]HandlerFunc
	staticPaths map[string]string // Map of URL paths to file directories
}

// NewRouter initializes a new Router
func NewRouter() *Router {
	return &Router{
		routes:      make(map[string]HandlerFunc),
		staticPaths: make(map[string]string),
	}
}

// AddRoute registers a new route with a handler
func (r *Router) AddRoute(method, path string, handler HandlerFunc) {
	key := fmt.Sprintf("%s:%s", method, path)
	r.routes[key] = handler
}

// AddStaticRoute registers a directory to serve static files
func (r *Router) AddStaticRoute(urlPath, dirPath string) {
	r.staticPaths[urlPath] = dirPath
}

// ServeHTTP matches requests to routes and serves files or calls handlers
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Check for static files first
	for urlPath, dirPath := range r.staticPaths {
		if req.URL.Path == urlPath || filepath.HasPrefix(req.URL.Path, urlPath+"/") {
			filePath := filepath.Join(dirPath, req.URL.Path[len(urlPath):])
			http.ServeFile(w, req, filePath)
			return
		}
	}

	// Check for dynamic routes
	key := fmt.Sprintf("%s:%s", req.Method, req.URL.Path)
	if handler, exists := r.routes[key]; exists {
		handler(w, req)
	} else {
		http.NotFound(w, req) // Return 404 if route is not found
	}
}

// RouteHandler adds routes to the router.
func (router *Router) RouteHandler() {
	router.AddRoute("POST", "/register", controllers.RegisterUserHandler)
	router.AddRoute("GET", "/register", controllers.RegisterUserHandler)
	router.AddRoute("POST", "/login", controllers.LoginHandler)
	router.AddRoute("POST", "/posts", controllers.PostsHandler)
	router.AddRoute("GET", "/", controllers.PostsHandler);
	router.AddRoute("GET", "/all_posts", controllers.GetAllPostsHandler)
	// Add other routes as needed
}
