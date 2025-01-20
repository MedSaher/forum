package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"forum/app/controllers"
	"forum/app/routers"
)

func main() {
	// Parse the static files:
	var err error
	controllers.Tmpl, err = template.ParseGlob("./app/views/*.html")
	if err != nil {
		log.Fatal(err)
	}
	// create a new instance of the router:
	router := routers.NewRouter()
	// Serve css static files:
	router.AddStaticRoute("/app/static/css", "./app/static/css")
	// Serve js static files:
	router.AddStaticRoute("/app/static/scripts", "./app/static/scripts")
	// Map the routs the specific handler
	router.RouteHandler()

	fmt.Println("run: http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", router))
}
