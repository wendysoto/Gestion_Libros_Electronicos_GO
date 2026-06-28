package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gestion_libros/services"

	"github.com/gorilla/mux"
)

type PrestamoHandler struct {
	prestamoService  *services.PrestamoService
	libroService     *services.LibroService
	usuarioService   *services.UsuarioService
	categoriaService *services.CategoriaService
	templates        *template.Template
}

func NuevoPrestamoHandler(
	prestamoService *services.PrestamoService,
	libroService *services.LibroService,
	usuarioService *services.UsuarioService,
	categoriaService *services.CategoriaService,
	tmpl *template.Template,
) *PrestamoHandler {
	return &PrestamoHandler{
		prestamoService:  prestamoService,
		libroService:     libroService,
		usuarioService:   usuarioService,
		categoriaService: categoriaService,
		templates:        tmpl,
	}
}

type PrestamoVista struct {
	ID              uint
	LibroTitulo     string
	UsuarioNombre   string
	FechaPrestamo   string
	FechaDevolucion string
	Estado          string
}

func (h *PrestamoHandler) ListarPrestamos(w http.ResponseWriter, r *http.Request) {
	prestamos, err := h.prestamoService.ListarPrestamos()
	if err != nil {
		log.Printf("Error al listar préstamos: %v", err)
		http.Error(w, "Error al obtener los préstamos", http.StatusInternalServerError)
		return
	}

	var prestamosVista []PrestamoVista
	for _, p := range prestamos {
		libroTitulo := "Desconocido"
		if libro, err := h.libroService.ObtenerLibro(p.GetLibroID()); err == nil {
			libroTitulo = libro.GetTitulo()
		}
		usuarioNombre := "Desconocido"
		if usr, err := h.usuarioService.ObtenerUsuario(p.GetUsuarioID()); err == nil {
			usuarioNombre = usr.GetNombre()
		}
		prestamosVista = append(prestamosVista, PrestamoVista{
			ID:              p.GetID(),
			LibroTitulo:     libroTitulo,
			UsuarioNombre:   usuarioNombre,
			FechaPrestamo:   p.GetFechaPrestamo().Format("2006-01-02"),
			FechaDevolucion: p.FechaDevolucionStr(),
			Estado:          string(p.GetEstado()),
		})
	}

	data := map[string]interface{}{
		"Titulo":    "Gestión de Préstamos",
		"Prestamos": prestamosVista,
	}
	h.templates.ExecuteTemplate(w, "prestamos.html", data)
}

func (h *PrestamoHandler) CrearPrestamoForm(w http.ResponseWriter, r *http.Request) {
	libros, _ := h.libroService.ListarLibros()
	usuarios, _ := h.usuarioService.ListarUsuarios()

	categorias, _ := h.categoriaService.ListarCategorias()
	catMap := make(map[uint]string)
	for _, c := range categorias {
		catMap[c.GetID()] = c.GetNombre()
	}

	type LibroVista struct {
		ID        uint
		Titulo    string
		Categoria string
	}
	var librosVista []LibroVista
	for _, l := range libros {
		nombreCat := ""
		if nombre, ok := catMap[l.GetCategoriaID()]; ok {
			nombreCat = nombre
		}
		librosVista = append(librosVista, LibroVista{
			ID:        l.GetID(),
			Titulo:    l.GetTitulo(),
			Categoria: nombreCat,
		})
	}

	data := map[string]interface{}{
		"Titulo":   "Registrar Nuevo Préstamo",
		"Libros":   librosVista,
		"Usuarios": usuarios,
	}
	h.templates.ExecuteTemplate(w, "prestamo_form.html", data)
}

// CrearPrestamo procesa el formulario. Valida límite de 3 préstamos por usuario.
func (h *PrestamoHandler) CrearPrestamo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	libroID, _ := strconv.ParseUint(r.FormValue("libro_id"), 10, 32)
	usuarioID, _ := strconv.ParseUint(r.FormValue("usuario_id"), 10, 32)

	_, err := h.prestamoService.CrearPrestamo(uint(libroID), uint(usuarioID))
	if err != nil {
		log.Printf("Error al crear préstamo: %v", err)

		libros, _ := h.libroService.ListarLibros()
		usuarios, _ := h.usuarioService.ListarUsuarios()
		categorias, _ := h.categoriaService.ListarCategorias()
		catMap := make(map[uint]string)
		for _, c := range categorias {
			catMap[c.GetID()] = c.GetNombre()
		}
		type LibroVista struct {
			ID        uint
			Titulo    string
			Categoria string
		}
		var librosVista []LibroVista
		for _, l := range libros {
			nombreCat := ""
			if nombre, ok := catMap[l.GetCategoriaID()]; ok {
				nombreCat = nombre
			}
			librosVista = append(librosVista, LibroVista{
				ID:        l.GetID(),
				Titulo:    l.GetTitulo(),
				Categoria: nombreCat,
			})
		}

		data := map[string]interface{}{
			"Titulo":   "Registrar Nuevo Préstamo",
			"Error":    err.Error(),
			"Libros":   librosVista,
			"Usuarios": usuarios,
		}
		h.templates.ExecuteTemplate(w, "prestamo_form.html", data)
		return
	}
	http.Redirect(w, r, "/prestamos", http.StatusSeeOther)
}

func (h *PrestamoHandler) DevolverPrestamo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	if err := h.prestamoService.DevolverPrestamo(uint(id)); err != nil {
		log.Printf("Error al devolver préstamo: %v", err)
		http.Error(w, "Error al registrar la devolución: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/prestamos", http.StatusSeeOther)
}

// HistorialPorUsuario muestra todos los préstamos de un usuario.
func (h *PrestamoHandler) HistorialPorUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	usuario, err := h.usuarioService.ObtenerUsuario(uint(id))
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	prestamos, err := h.prestamoService.ObtenerPorUsuario(uint(id))
	if err != nil {
		log.Printf("Error al obtener historial: %v", err)
		http.Error(w, "Error al obtener historial", http.StatusInternalServerError)
		return
	}

	var prestamosVista []PrestamoVista
	for _, p := range prestamos {
		libroTitulo := "Desconocido"
		if libro, err := h.libroService.ObtenerLibro(p.GetLibroID()); err == nil {
			libroTitulo = libro.GetTitulo()
		}
		prestamosVista = append(prestamosVista, PrestamoVista{
			ID:              p.GetID(),
			LibroTitulo:     libroTitulo,
			UsuarioNombre:   usuario.GetNombre(),
			FechaPrestamo:   p.GetFechaPrestamo().Format("2006-01-02"),
			FechaDevolucion: p.FechaDevolucionStr(),
			Estado:          string(p.GetEstado()),
		})
	}

	data := map[string]interface{}{
		"Titulo":    "Historial de Préstamos - " + usuario.GetNombre(),
		"Prestamos": prestamosVista,
		"Usuario":   usuario,
	}
	h.templates.ExecuteTemplate(w, "historial.html", data)
}
