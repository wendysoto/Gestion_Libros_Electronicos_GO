// Package handlers contiene los controladores HTTP que manejan las peticiones
// del frontend y responden con las vistas HTML renderizadas.
package handlers

import (
	"html/template"
	"log"
	"net/http"

	"gestion_libros/services"
)

// LibroHandler gestiona las peticiones HTTP relacionadas con libros.
type LibroHandler struct {
	service   *services.LibroService
	templates *template.Template
}

// NuevoLibroHandler crea una instancia del handler con sus dependencias.
func NuevoLibroHandler(service *services.LibroService, tmpl *template.Template) *LibroHandler {
	return &LibroHandler{service: service, templates: tmpl}
}

// ListarLibros consulta todos los libros desde la base de datos
// y renderiza la vista HTML con el catálogo completo.
func (h *LibroHandler) ListarLibros(w http.ResponseWriter, r *http.Request) {
	libros, err := h.service.ListarLibros()
	if err != nil {
		log.Printf("Error al listar libros: %v", err)
		http.Error(w, "Error al obtener los libros", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Titulo": "Catálogo de Libros",
		"Libros": libros,
	}

	if err := h.templates.ExecuteTemplate(w, "libros.html", data); err != nil {
		log.Printf("Error al renderizar template: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}
