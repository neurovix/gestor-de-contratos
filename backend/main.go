package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	var err error

	if err = godotenv.Load(".env"); err != nil {
		log.Println(err.Error())
	}

	var (
		routerPort string = os.Getenv("ROUTER_PORT")
	)

	var mux = http.NewServeMux()

	if err = http.ListenAndServe(routerPort, mux); err != nil {
		log.Println(err.Error())
	}
}
