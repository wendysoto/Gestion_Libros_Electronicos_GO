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

type UsuarioHandler struct {
	service   *services.UsuarioService
	templates *template.Template
}

func NuevoUsuarioHandler(service *services.UsuarioService, tmpl *template.Template) *UsuarioHandler {
	return &UsuarioHandler{service: service, templates: tmpl}
}

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
	h.templates.ExecuteTemplate(w, "usuarios.html", data)
}

func (h *UsuarioHandler) CrearUsuarioForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Titulo": "Registrar Nuevo Usuario",
	}
	h.templates.ExecuteTemplate(w, "usuario_form.html", data)
}

func (h *UsuarioHandler) CrearUsuario(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nombre := r.FormValue("nombre")
	email := r.FormValue("email")
	tipo := models.TipoUsuario(r.FormValue("tipo"))

	_, err := h.service.CrearUsuario(nombre, email, tipo)
	if err != nil {
		log.Printf("Error al crear usuario: %v", err)
		data := map[string]interface{}{
			"Titulo": "Registrar Nuevo Usuario",
			"Error":  err.Error(),
		}
		h.templates.ExecuteTemplate(w, "usuario_form.html", data)
		return
	}
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}

// EditarUsuarioForm muestra el formulario con los datos actuales del usuario.
func (h *UsuarioHandler) EditarUsuarioForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	usuario, err := h.service.ObtenerUsuario(uint(id))
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}
	data := map[string]interface{}{
		"Titulo":  "Editar Usuario",
		"Usuario": usuario,
	}
	h.templates.ExecuteTemplate(w, "usuario_editar.html", data)
}

// EditarUsuario procesa la edición de un usuario.
func (h *UsuarioHandler) EditarUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	usuario, err := h.service.ObtenerUsuario(uint(id))
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	r.ParseForm()
	usuario.SetNombre(r.FormValue("nombre"))
	usuario.SetEmail(r.FormValue("email"))
	usuario.SetTipo(models.TipoUsuario(r.FormValue("tipo")))

	if err := h.service.ActualizarUsuario(usuario); err != nil {
		log.Printf("Error al editar usuario: %v", err)
		data := map[string]interface{}{
			"Titulo":  "Editar Usuario",
			"Usuario": usuario,
			"Error":   err.Error(),
		}
		h.templates.ExecuteTemplate(w, "usuario_editar.html", data)
		return
	}
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}

// EliminarUsuario elimina un usuario solo si no tiene préstamos activos.
func (h *UsuarioHandler) EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	if err := h.service.EliminarUsuario(uint(id)); err != nil {
		log.Printf("Error al eliminar usuario: %v", err)
		// Redirigir con mensaje de error en query param
		http.Redirect(w, r, "/usuarios?error="+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}
