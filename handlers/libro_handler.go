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

type LibroHandler struct {
	libroService     *services.LibroService
	categoriaService *services.CategoriaService
	templates        *template.Template
}

func NuevoLibroHandler(libroService *services.LibroService, categoriaService *services.CategoriaService, tmpl *template.Template) *LibroHandler {
	return &LibroHandler{
		libroService:     libroService,
		categoriaService: categoriaService,
		templates:        tmpl,
	}
}

type LibroVista struct {
	ID          uint
	Titulo      string
	Autor       string
	Categoria   string
	CategoriaID uint
	ISBN        string
	Formato     string
}

func (h *LibroHandler) construirLibrosVista(libros []*models.Libro, catMap map[uint]string) []LibroVista {
	vista := make([]LibroVista, len(libros))
	for i, l := range libros {
		nombreCat := "Sin categoría"
		if nombre, ok := catMap[l.GetCategoriaID()]; ok {
			nombreCat = nombre
		}
		vista[i] = LibroVista{
			ID:          l.GetID(),
			Titulo:      l.GetTitulo(),
			Autor:       l.GetAutor(),
			Categoria:   nombreCat,
			CategoriaID: l.GetCategoriaID(),
			ISBN:        l.GetISBN(),
			Formato:     string(l.GetFormato()),
		}
	}
	return vista
}

func (h *LibroHandler) obtenerCatMap() map[uint]string {
	catMap := make(map[uint]string)
	categorias, err := h.categoriaService.ListarCategorias()
	if err == nil {
		for _, c := range categorias {
			catMap[c.GetID()] = c.GetNombre()
		}
	}
	return catMap
}

func (h *LibroHandler) ListarLibros(w http.ResponseWriter, r *http.Request) {
	libros, err := h.libroService.ListarLibros()
	if err != nil {
		log.Printf("Error al listar libros: %v", err)
		http.Error(w, "Error al obtener los libros", http.StatusInternalServerError)
		return
	}
	categorias, _ := h.categoriaService.ListarCategorias()
	data := map[string]interface{}{
		"Titulo":     "Catálogo de Libros",
		"Libros":     h.construirLibrosVista(libros, h.obtenerCatMap()),
		"Categorias": categorias,
	}
	h.templates.ExecuteTemplate(w, "libros.html", data)
}

func (h *LibroHandler) CrearLibroForm(w http.ResponseWriter, r *http.Request) {
	categorias, _ := h.categoriaService.ListarCategorias()
	data := map[string]interface{}{
		"Titulo":     "Registrar Nuevo Libro",
		"Categorias": categorias,
	}
	h.templates.ExecuteTemplate(w, "libro_form.html", data)
}

func (h *LibroHandler) CrearLibro(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	titulo := r.FormValue("titulo")
	autor := r.FormValue("autor")
	categoriaIDStr := r.FormValue("categoria_id")
	isbn := r.FormValue("isbn")
	formato := models.FormatoLibro(r.FormValue("formato"))

	categoriaID, err := strconv.ParseUint(categoriaIDStr, 10, 32)
	if err != nil {
		categorias, _ := h.categoriaService.ListarCategorias()
		data := map[string]interface{}{
			"Titulo":     "Registrar Nuevo Libro",
			"Error":      "Debe seleccionar una categoría válida",
			"Categorias": categorias,
		}
		h.templates.ExecuteTemplate(w, "libro_form.html", data)
		return
	}

	_, err = h.libroService.CrearLibro(titulo, autor, uint(categoriaID), isbn, formato)
	if err != nil {
		log.Printf("Error al crear libro: %v", err)
		categorias, _ := h.categoriaService.ListarCategorias()
		data := map[string]interface{}{
			"Titulo":     "Registrar Nuevo Libro",
			"Error":      err.Error(),
			"Categorias": categorias,
		}
		h.templates.ExecuteTemplate(w, "libro_form.html", data)
		return
	}
	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}

// EditarLibroForm muestra el formulario con los datos actuales del libro.
func (h *LibroHandler) EditarLibroForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	libro, err := h.libroService.ObtenerLibro(uint(id))
	if err != nil {
		http.Error(w, "Libro no encontrado", http.StatusNotFound)
		return
	}
	categorias, _ := h.categoriaService.ListarCategorias()
	data := map[string]interface{}{
		"Titulo":     "Editar Libro",
		"Libro":      libro,
		"Categorias": categorias,
	}
	h.templates.ExecuteTemplate(w, "libro_editar.html", data)
}

// EditarLibro procesa la edición de un libro.
func (h *LibroHandler) EditarLibro(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	libro, err := h.libroService.ObtenerLibro(uint(id))
	if err != nil {
		http.Error(w, "Libro no encontrado", http.StatusNotFound)
		return
	}

	r.ParseForm()
	libro.SetTitulo(r.FormValue("titulo"))
	libro.SetAutor(r.FormValue("autor"))
	catID, _ := strconv.ParseUint(r.FormValue("categoria_id"), 10, 32)
	libro.SetCategoriaID(uint(catID))
	libro.SetISBN(r.FormValue("isbn"))
	libro.SetFormato(models.FormatoLibro(r.FormValue("formato")))

	if err := h.libroService.ActualizarLibro(libro); err != nil {
		log.Printf("Error al editar libro: %v", err)
		categorias, _ := h.categoriaService.ListarCategorias()
		data := map[string]interface{}{
			"Titulo":     "Editar Libro",
			"Libro":      libro,
			"Error":      err.Error(),
			"Categorias": categorias,
		}
		h.templates.ExecuteTemplate(w, "libro_editar.html", data)
		return
	}
	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}

func (h *LibroHandler) EliminarLibro(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	if err := h.libroService.EliminarLibro(uint(id)); err != nil {
		log.Printf("Error al eliminar libro: %v", err)
		http.Error(w, "Error al eliminar el libro", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}

func (h *LibroHandler) BuscarLibrosPorCategoria(w http.ResponseWriter, r *http.Request) {
	categoriaIDStr := r.URL.Query().Get("categoria_id")
	categorias, _ := h.categoriaService.ListarCategorias()
	catMap := h.obtenerCatMap()

	var libros []*models.Libro
	var err error
	if categoriaIDStr == "" || categoriaIDStr == "0" {
		libros, err = h.libroService.ListarLibros()
	} else {
		catID, _ := strconv.ParseUint(categoriaIDStr, 10, 32)
		libros, err = h.libroService.BuscarPorCategoria(uint(catID))
	}
	if err != nil {
		log.Printf("Error al buscar libros: %v", err)
		http.Error(w, "Error al buscar libros", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Titulo":                "Catálogo de Libros",
		"Libros":                h.construirLibrosVista(libros, catMap),
		"Categorias":            categorias,
		"CategoriaSeleccionada": categoriaIDStr,
	}
	h.templates.ExecuteTemplate(w, "libros.html", data)
}
