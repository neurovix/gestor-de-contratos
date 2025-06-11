package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/neurovix/tramites/backend/src/database"
)

func DeleteTramite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var idTramite = r.URL.Query().Get("id-tramite")

	if _, err := database.DB.Exec(context.Background(), "delete from tramites where id_tramite = $1", &idTramite); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status_code": http.StatusOK,
		"message":     "tramite borrado exitosamente",
	})
}
