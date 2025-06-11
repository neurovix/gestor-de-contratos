package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/neurovix/tramites/backend/src/database"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	var userID int
	row := database.DB.QueryRow(context.Background(), "SELECT id_usuario, password FROM usuarios WHERE email = $1", creds.Email)
	if err := row.Scan(&userID, &hashedPassword); err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
		return
	}

	// ✅ Cookie corregida para CORS
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    strconv.Itoa(userID),
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false, // ⭐ CAMBIO CLAVE: false para permitir acceso desde JavaScript
		Secure:   false, // Solo true si usas HTTPS
		Path:     "/",
		Domain:   "",                   // ⭐ CAMBIO: vacío es más compatible para localhost
		SameSite: http.SameSiteLaxMode, // ⭐ AÑADIDO: importante para CORS
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Inicio de sesión exitoso"})
}
