package models

type VerificadorInput struct {
	Nombre     string `json:"nombre"`
	Planta     string `json:"planta"`
	Orden      int    `json:"orden"`
	Verificado bool   `json:"verificado"`
}

type TrámiteRequest struct {
	NombreTramite string             `json:"nombre_tramite"`
	ArchivoPDFURL string             `json:"archivo_pdf_url"`
	CreadorID     int                `json:"creador_id"`
	Verificadores []VerificadorInput `json:"verificadores"`
}
