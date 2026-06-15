// Package handlers contiene los controladores HTTP que manejan las peticiones
// del frontend y responden con las vistas HTML renderizadas.
package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gestion_libros/models"
	"gestion_libros/services"

	"github.com/gorilla/mux"
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

// ListarLibros muestra la página principal con todos los libros.
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

// CrearLibroForm muestra el formulario para crear un nuevo libro.
func (h *LibroHandler) CrearLibroForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Titulo": "Registrar Nuevo Libro",
	}
	if err := h.templates.ExecuteTemplate(w, "libro_form.html", data); err != nil {
		log.Printf("Error al renderizar template: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// CrearLibro procesa el formulario y crea un nuevo libro.
func (h *LibroHandler) CrearLibro(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	titulo := r.FormValue("titulo")
	autor := r.FormValue("autor")
	categoria := r.FormValue("categoria")
	isbn := r.FormValue("isbn")
	formato := models.FormatoLibro(r.FormValue("formato"))

	_, err := h.service.CrearLibro(titulo, autor, categoria, isbn, formato)
	if err != nil {
		log.Printf("Error al crear libro: %v", err)
		data := map[string]interface{}{
			"Titulo": "Registrar Nuevo Libro",
			"Error":  err.Error(),
		}
		h.templates.ExecuteTemplate(w, "libro_form.html", data)
		return
	}

	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}

// EliminarLibro elimina un libro por su ID.
func (h *LibroHandler) EliminarLibro(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.EliminarLibro(uint(id)); err != nil {
		log.Printf("Error al eliminar libro: %v", err)
		http.Error(w, "Error al eliminar el libro", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}
