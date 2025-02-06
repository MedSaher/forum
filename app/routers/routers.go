package routers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"forum/app/controllers"
	"forum/app/models"
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
		http.NotFound(w, req) // Return 404 if route is not found
	}
}

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
	router.AddRoute("POST", "/vote", controllers.CreateVote)
}

// Add a middleware for static files:
func (router *Router) StaticMiddleWare() {
	// Serve css static files:
	router.AddStaticRoute("/app/static/css", "./app/static/css")
	// Serve js static files:
	router.AddStaticRoute("/app/static/scripts", "./app/static/scripts")
	// serve pictures:
	router.AddStaticRoute("/app/uploads", "./app/uploads")
}

// Create a session middleware in case of the abcence of login the program will force the user to log in:
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session UUID from the cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Validate the session from the database
		userID, err := models.GetSessionByUUID(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Add the userID to the request context
		ctx := context.WithValue(r.Context(), "UserID", userID)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
