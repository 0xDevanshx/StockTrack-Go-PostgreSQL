package main

import (
	"fmt"
	"go-postgres/router"
	"log"
	"net/http"
	_ "go-postgres/database" // Import the database package so that init() opens the pool.
	
)

func main(){
	// At this point, database.init() has already:
	//   1) loaded .env
	//   2) set up the global pool in database.DB

	r := router.Router()
	fmt.Println("Starting server on port 6969...")

	log.Fatal(http.ListenAndServe(":6969",r))
}