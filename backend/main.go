package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/neurovix/tramites/backend/src/auth"
	"github.com/neurovix/tramites/backend/src/database"
	"github.com/neurovix/tramites/backend/src/handlers"
)

func main() {
	var err error

	if err = database.InitDB(); err != nil {
		log.Println(err.Error())
	}

	log.Println("Conexion exitosa a la base de datos")

	defer database.DB.Close()

	if err = godotenv.Load(".env"); err != nil {
		log.Println(err.Error())
	}

	var (
		routerPort string = os.Getenv("ROUTER_PORT")
	)

	var mux = http.NewServeMux()

	mux.HandleFunc("POST /api/login", auth.LoginHandler)
	mux.HandleFunc("POST /api/register", auth.RegisterHandler)
	mux.HandleFunc("POST /api/logout", auth.LogoutHandler)
	// mux.HandleFunc("GET /api/plantas", handlers.GetPlantas)
	mux.HandleFunc("POST /api/nuevo_tramite", handlers.NuevoTramite)
	mux.HandleFunc("POST /api/upload_file", handlers.UploadFile)

	var router = WithCors(mux)

	log.Println("Servidor corriendo en http://localhost:8080/api")

	if err = http.ListenAndServe(routerPort, router); err != nil {
		log.Println(err.Error())
	}
}

func WithCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "http://localhost:4321" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
