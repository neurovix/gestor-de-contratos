package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/neurovix/tramites/backend/src/database"
)

func GetEmployeesBasedOnPlant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var idPlanta = r.URL.Query().Get("id_planta")

	rows, err := database.DB.Query(context.Background(), "select id_usuario, nombre from usuarios where id_planta = $1", idPlanta)

	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, "error al obtener el id_planta")
		return
	}

	defer rows.Close()

	type Usuario struct {
		Id     string `json:"id_usuario"`
		Nombre string `json:"nombre_usuario"`
	}

	var usuarios = []Usuario{}

	for rows.Next() {
		var usuario = new(Usuario)

		if err := rows.Scan(&usuario.Id, &usuario.Nombre); err != nil {
			sendJSONError(w, http.StatusConflict, "error al procesar los usuarios de la planta con id "+idPlanta)
			return
		}

		usuarios = append(usuarios, *usuario)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status_code": http.StatusOK,
		"body":        usuarios,
	})
}

type Seguimiento struct {
	Orden      int    `json:"orden"`
	Nombre     string `json:"nombre"`
	Planta     string `json:"planta"`
	Verificado bool   `json:"verificado"`
}

type Tramite struct {
	NoContrato       string        `json:"no_contrato"`
	CreadoPor        string        `json:"creado_por"`
	CreadoEn         string        `json:"creado_en"`
	UrlArchivo       string        `json:"url_archivo"`
	ListaSeguimiento []Seguimiento `json:"lista_seguimiento"`
	PuedeAprobar     bool          `json:"puede_aprobar"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func GetTramite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	idTramite := r.URL.Query().Get("id-tramite")
	if idTramite == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Falta el parámetro 'id-tramite'"})
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Usuario no autenticado"})
		return
	}
	var idUsuario = cookie.Value

	var (
		noContrato          string
		creadoPor           string
		creadoEn            string
		urlArchivo          string
		listaSeguimientoRaw []byte
	)

	query := `
		SELECT 
			t.no_contrato,
			CONCAT(p.nombre, ', ', u.nombre) AS creado_por,
			TO_CHAR(t.fecha_creacion, 'YYYY-MM-DD HH24:MI:SS') AS creado_en,
			t.archivo_pdf_url AS url_archivo,
			(
				SELECT json_agg(json_build_object(
					'orden', v.orden,
					'nombre', uv.nombre,
					'planta', pv.nombre,
					'verificado', v.verificado
				) ORDER BY v.orden)
				FROM verificaciones v
				JOIN usuarios uv ON v.id_verificador = uv.id_usuario
				JOIN plantas pv ON uv.id_planta = pv.id_planta
				WHERE v.id_tramite = t.id_tramite
			) AS lista_de_seguimiento
		FROM tramites t
		JOIN usuarios u ON t.creador_id = u.id_usuario
		JOIN plantas p ON u.id_planta = p.id_planta
		WHERE t.id_tramite = $1
	`

	row := database.DB.QueryRow(context.Background(), query, idTramite)
	err = row.Scan(&noContrato, &creadoPor, &creadoEn, &urlArchivo, &listaSeguimientoRaw)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No se encontró el trámite o error en la consulta: " + err.Error()})
		return
	}

	var puedeAprobar bool

	querySiguiente := `
	SELECT id_verificador
	FROM verificaciones
	WHERE id_tramite = $1 AND verificado = false
	ORDER BY orden ASC
	LIMIT 1
`
	var idSiguienteVerificador int
	err = database.DB.QueryRow(context.Background(), querySiguiente, idTramite).Scan(&idSiguienteVerificador)
	if err == nil && strconv.Itoa(idSiguienteVerificador) == idUsuario {
		puedeAprobar = true
	}

	var listaSeguimiento []Seguimiento
	if err := json.Unmarshal(listaSeguimientoRaw, &listaSeguimiento); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error al procesar la lista de seguimiento: " + err.Error()})
		return
	}

	tramite := Tramite{
		NoContrato:       noContrato,
		CreadoPor:        creadoPor,
		CreadoEn:         creadoEn,
		UrlArchivo:       urlArchivo,
		ListaSeguimiento: listaSeguimiento,
		PuedeAprobar:     puedeAprobar,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tramite)
}

func GetTramites(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Error obteniendo cookie session_id: %v", err)
		log.Println("Headers recibidos:", r.Header)
		sendJSONError(w, http.StatusBadRequest, "error al obtener la cookie: "+err.Error())
		return
	}

	id_usuario, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Printf("Error parseando cookie value '%s': %v", cookie.Value, err)
		sendJSONError(w, http.StatusInternalServerError, "error al parsear la cookie: "+err.Error())
		return
	}

	rows, err := database.DB.Query(context.Background(),
		"SELECT t.id_tramite, t.no_contrato, t.nombre_tramite FROM tramites t JOIN verificaciones v ON v.id_tramite = t.id_tramite WHERE v.id_verificador = $1",
		id_usuario)

	if err != nil {
		log.Printf("Error en query de base de datos: %v", err)
		sendJSONError(w, http.StatusInternalServerError, "error interno: "+err.Error())
		return
	}
	defer rows.Close()

	type TramiteHome struct {
		Id         int    `json:"id_tramite"`
		NoContrato string `json:"no_contrato"`
		Nombre     string `json:"nombre_tramite"`
	}

	var tramites []TramiteHome

	for rows.Next() {
		var t = new(TramiteHome)
		if err := rows.Scan(&t.Id, &t.NoContrato, &t.Nombre); err != nil {
			log.Printf("Error escaneando fila: %v", err)
			sendJSONError(w, http.StatusInternalServerError, "error escaneando trámite: "+err.Error())
			return
		}
		tramites = append(tramites, *t)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status_code": http.StatusOK,
		"body":        tramites,
	})
}
