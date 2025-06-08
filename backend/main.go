package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/neurovix/tramites/backend/src/database"
)

func main() {
	var err error

	if _, err = database.ConnDB(); err != nil {
		log.Println(err.Error())
	}

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
