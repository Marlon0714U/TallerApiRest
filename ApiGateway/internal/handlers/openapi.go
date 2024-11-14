package handlers

import (
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

// OpenAPIDocHandler sirve el archivo openapi.yml como JSON en /docs
func OpenAPIDocHandler(w http.ResponseWriter, r *http.Request) {
	// Leer el archivo openapi.yml
	data, err := os.ReadFile("openapi.yml")
	if err != nil {
		http.Error(w, "No se pudo leer el archivo OpenAPI", http.StatusInternalServerError)
		return
	}

	// Parsear el YAML a un objeto OpenAPI
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		http.Error(w, "No se pudo parsear el archivo OpenAPI", http.StatusInternalServerError)
		return
	}

	// Serializar a JSON y responder
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := doc.MarshalJSON()
	if err != nil {
		http.Error(w, "No se pudo convertir el archivo OpenAPI a JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
