package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/neurovix/tramites/backend/src/database"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var user struct {
		Nombre   string `json:"nombre"`
		Email    string `json:"email"`
		Password string `json:"password"`
		IDPlanta int    `json:"id_planta"`
		Cargo    string `json:"cargo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al encriptar la contraseña", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	_, err = database.DB.Exec(
		context.Background(),
		"INSERT INTO usuarios (nombre, email, password, id_planta, cargo) VALUES ($1, $2, $3, $4, $5)",
		user.Nombre, user.Email, string(hash), user.IDPlanta, user.Cargo,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error al registrar el usuario",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario registrado exitosamente"})
}
