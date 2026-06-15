package handlers

import (
	"html/template"
	"log"
	"net/http"

	"gestion_libros/services"
)

// UsuarioHandler gestiona las peticiones HTTP de usuarios.
type UsuarioHandler struct {
	service   *services.UsuarioService
	templates *template.Template
}

// NuevoUsuarioHandler crea una instancia del handler de usuarios.
func NuevoUsuarioHandler(service *services.UsuarioService, tmpl *template.Template) *UsuarioHandler {
	return &UsuarioHandler{service: service, templates: tmpl}
}

// ListarUsuarios consulta todos los usuarios desde la base de datos
// y renderiza la vista HTML con la lista de usuarios registrados.
func (h *UsuarioHandler) ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.service.ListarUsuarios()
	if err != nil {
		log.Printf("Error al listar usuarios: %v", err)
		http.Error(w, "Error al obtener los usuarios", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Titulo":   "Gestión de Usuarios",
		"Usuarios": usuarios,
	}

	if err := h.templates.ExecuteTemplate(w, "usuarios.html", data); err != nil {
		log.Printf("Error al renderizar template: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}
