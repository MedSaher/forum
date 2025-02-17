package routers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

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
		// Check if the request URL path starts with the static path
		if strings.HasPrefix(req.URL.Path, urlPath+"/") || req.URL.Path == urlPath {
			// Safely join the directory path with the file path
			filePath := filepath.Join(dirPath, req.URL.Path[len(urlPath):])
			// Normalize the file path
			filePath = filepath.Clean(filePath)
			// Ensure it's not a directory traversal attack
			if strings.Contains(filePath, "..") {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			// Serve the static file
			http.ServeFile(w, req, filePath)
			return
		}
	}
	// Check for dynamic routes
	key := fmt.Sprintf("%s:%s", req.Method, req.URL.Path)
	if handler, exists := r.routes[key]; exists {
		handler(w, req)
	} else {
		// http.NotFound(w, req) // Return 404 if route is not found
		controllers.ErrorHandler(w, req, 404)
	}
}

// RouteHandler adds routes to the router.
// RouteHandler adds routes to the router.
func (router *Router) MiddleWare() {
	router.AddRoute("POST", "/register", controllers.RegisterUserHandler)
	router.AddRoute("GET", "/register", controllers.RegisterUserHandler)
	router.AddRoute("POST", "/login", controllers.LoginHandler)
	router.AddRoute("GET", "/login", controllers.LoginHandler)
	router.AddRoute("GET", "/", controllers.PostsHandler)
	router.AddRoute("GET", "/all_posts", controllers.GetAllPostsHandler)
	router.AddRoute("GET", "/all_categories", controllers.GetAllCategories)
	router.AddRoute("GET", "/profile", controllers.LogedInUser)
	router.AddRoute("POST", "/logout", controllers.Logout)
	router.AddRoute("GET", "/add_post", controllers.AddPost)
	router.AddRoute("POST", "/add_post", controllers.AddPost)
	router.AddRoute("POST", "/vote_for_post", controllers.VoteForPost)
	router.AddRoute("GET", "/liked", controllers.GetLikedPosts)
	router.AddRoute("GET", "/owned", controllers.GetOwnedPosts)
	router.AddRoute("GET", "/get_comments", controllers.GetComments)
	router.AddRoute("POST", "/post_comment", controllers.CreateComment)
	router.AddRoute("POST", "/vote_for_comment", controllers.VoteForComment)
}

// Add a middleware for static files:
func (router *Router) StaticMiddleWare() {
	// Serve css static files:
	router.AddStaticRoute("/app/static/css", "./app/static/css")
	// Serve js static files:
	router.AddStaticRoute("/app/static/editjs", "./app/static/editjs")
	// serve pictures:
	router.AddStaticRoute("/app/uploads", "./app/uploads")
}
