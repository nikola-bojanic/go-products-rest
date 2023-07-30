package main

import (
	"fmt"
	"go-products-rest/router"
	"log"
	"net/http"
)

func main() {
	router := router.Router()
	fmt.Println("starting server on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
