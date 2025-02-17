package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"forum/app/config"
	"forum/app/controllers"
	"forum/app/models"
	"forum/app/routers"
)

func main() {
	var err error
	// Check if the database file exists
	if _, err = os.Stat("forum.db"); os.IsNotExist(err) {
		// Create an empty file
		file, err := os.Create("forum.db")
		if err != nil {
			log.Fatal(err)
		}
		// Instantiate the schema:
		err = config.CreateSchema()
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}
	closeDB := make(chan os.Signal, 1)
	go closeDBFunc(closeDB)
	// Parse the static files:
	controllers.Tmpl, err = template.ParseGlob("./app/views/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// create a new instance of the router:
	router := routers.NewRouter()
	// handle static files:
	router.StaticMiddleWare()
	// Map the routs the specific handler
	router.MiddleWare()

	fmt.Println("run: http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func closeDBFunc(closeChan chan os.Signal) {
	signal.Notify(closeChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-closeChan
	models.DeleteAllSessions()
	os.Exit(0)
}
