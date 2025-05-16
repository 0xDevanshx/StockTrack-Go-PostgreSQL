package main

import (
	"fmt"
	"go-postgres/router"
	"log"
	"net/http"
	
)

func main(){
	r := router.Router()
	fmt.Println("Starting server on port 6969...")

	log.Fatal(http.ListenAndServe(":6969",r))
}