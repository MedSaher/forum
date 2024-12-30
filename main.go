package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/app/routers"
)

func main() {
	// declare a new instance of the router:
	router := routers.NewRouter()
	// Map the routs the specific handler
	router.RouteHandler()

	// Start the server
	fmt.Println("run: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
