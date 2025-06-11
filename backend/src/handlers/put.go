package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/neurovix/tramites/backend/src/database"
)

func UpdateTramite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	cookie, err := r.Cookie("session_id")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	var idUsuario = cookie.Value

	var idTramite = r.URL.Query().Get("id-tramite")

	if _, err = database.DB.Exec(context.Background(), "update verificaciones set verificado = true where id_tramite = $1 and id_verificador = $2", &idTramite, &idUsuario); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Aprobación registrada exitosamente",
	})
}
