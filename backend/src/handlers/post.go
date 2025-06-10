package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/neurovix/tramites/backend/src/database"
)

type VerificadorInput struct {
	Nombre     string `json:"nombre"`
	Planta     string `json:"planta"`
	Orden      int    `json:"orden"`
	Verificado bool   `json:"verificado"`
}

type TramiteRequest struct {
	NombreTramite string             `json:"nombre_tramite"`
	ArchivoPDFURL string             `json:"archivo_pdf_url"`
	CreadorID     int                `json:"creador_id"`
	Verificadores []VerificadorInput `json:"verificadores"`
}

func NuevoTramite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var datos TramiteRequest
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		sendJSONError(w, http.StatusBadRequest, "Error al decodificar JSON: "+err.Error())
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		sendJSONError(w, http.StatusUnauthorized, "Error al obtener la cookie: "+err.Error())
		return
	}

	if datos.CreadorID, err = strconv.Atoi(cookie.Value); err != nil {
		sendJSONError(w, http.StatusBadRequest, "Error al convertir el session_id en int: "+err.Error())
		return
	}
	// Verificadores con ID
	type VerificadorFinal struct {
		IDVerificador int  `json:"id_verificador"`
		Orden         int  `json:"orden"`
		Verificado    bool `json:"verificado"`
	}

	var verificadoresFinales []VerificadorFinal

	for _, v := range datos.Verificadores {
		var idVerificador int
		err := database.DB.QueryRow(context.Background(), `
			SELECT u.id_usuario
			FROM usuarios u
			JOIN plantas p ON u.id_planta = p.id_planta
			WHERE u.nombre = $1 AND p.nombre = $2
		`, v.Nombre, v.Planta).Scan(&idVerificador)

		if err != nil {
			sendJSONError(w, http.StatusBadRequest, "No se encontró el verificador: "+v.Nombre+" en "+v.Planta+" → "+err.Error())
			return
		}

		verificadoresFinales = append(verificadoresFinales, VerificadorFinal{
			IDVerificador: idVerificador,
			Orden:         v.Orden,
			Verificado:    v.Verificado,
		})
	}

	log.Println(verificadoresFinales)

	verificadoresJSON, err := json.Marshal(verificadoresFinales)
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, "Error al convertir verificadores a JSON: "+err.Error())
		return
	}

	var idTramite int
	err = database.DB.QueryRow(context.Background(), `
		SELECT crear_tramite_con_verificadores($1, $2, $3, $4)
	`, datos.NombreTramite, datos.ArchivoPDFURL, datos.CreadorID, string(verificadoresJSON)).Scan(&idTramite)

	log.Println(idTramite)

	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, "Error al insertar en la base de datos: "+err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"mensaje":    "Trámite creado correctamente",
		"id_tramite": idTramite,
	})
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	file, header, err := r.FormFile("pdf")
	if err != nil {
		sendJSONError(w, http.StatusBadRequest, "Error al recibir el archivo: "+err.Error())
		return
	}
	defer file.Close()

	path := "./uploads/" + header.Filename
	out, err := os.Create(path)
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, "Error al guardar el archivo: "+err.Error())
		return
	}
	defer out.Close()

	io.Copy(out, file)
	json.NewEncoder(w).Encode(map[string]string{"url": "/uploads/" + header.Filename})
}

func sendJSONError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
