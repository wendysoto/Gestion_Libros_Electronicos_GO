// Package handlers - API REST en formato JSON.
// Endpoint para listar libros por categoría retornando JSON.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gestion_libros/models"
	"gestion_libros/services"
)

// APIHandler gestiona los endpoints REST que retornan JSON.
type APIHandler struct {
	libroService     *services.LibroService
	categoriaService *services.CategoriaService
	usuarioService   *services.UsuarioService
	prestamoService  *services.PrestamoService
}

// NuevoAPIHandler crea una instancia del handler de API.
func NuevoAPIHandler(
	libroService *services.LibroService,
	categoriaService *services.CategoriaService,
	usuarioService *services.UsuarioService,
	prestamoService *services.PrestamoService,
) *APIHandler {
	return &APIHandler{
		libroService:     libroService,
		categoriaService: categoriaService,
		usuarioService:   usuarioService,
		prestamoService:  prestamoService,
	}
}

// LibroJSON es la estructura de respuesta JSON para el endpoint de libros.
type LibroJSON struct {
	Titulo     string `json:"titulo"`
	Autor      string `json:"autor"`
	Categoria  string `json:"categoria"`
	Disponible string `json:"disponible"`
}

type PrestamoPorUsuarioRequest struct {
	NombreUsuario string `json:"nombre_usuario"`
}

type PrestamoJSON struct {
	ID              uint   `json:"id"`
	Libro           string `json:"libro"`
	Usuario         string `json:"usuario"`
	FechaPrestamo   string `json:"fecha_prestamo"`
	FechaDevolucion string `json:"fecha_devolucion"`
	Estado          string `json:"estado"`
}

// ListarLibrosPorCategoria retorna los libros filtrados por categoría en formato JSON.
// Endpoint: GET /api/libros?categoria_id=1
// Respuesta: [{"titulo":"...", "autor":"...", "categoria":"...", "disponible":"SI"}]
func (h *APIHandler) ListarLibrosPorCategoria(w http.ResponseWriter, r *http.Request) {
	categoriaIDStr := r.URL.Query().Get("categoria_id")

	// Obtener mapa de categorías
	categorias, err := h.categoriaService.ListarCategorias()
	if err != nil {
		log.Printf("Error al listar categorías: %v", err)
		http.Error(w, `{"error":"Error al obtener categorías"}`, http.StatusInternalServerError)
		return
	}
	catMap := make(map[uint]string)
	for _, c := range categorias {
		catMap[c.GetID()] = c.GetNombre()
	}

	// Obtener libros (todos o filtrados por categoría)
	var librosJSON []LibroJSON
	if categoriaIDStr == "" || categoriaIDStr == "0" {
		libros, err := h.libroService.ListarLibros()
		if err != nil {
			log.Printf("Error al listar libros: %v", err)
			http.Error(w, `{"error":"Error al obtener libros"}`, http.StatusInternalServerError)
			return
		}
		for _, l := range libros {
			disponible := "SI"
			if !l.GetDisponible() {
				disponible = "NO"
			}
			nombreCat := "Sin categoría"
			if nombre, ok := catMap[l.GetCategoriaID()]; ok {
				nombreCat = nombre
			}
			librosJSON = append(librosJSON, LibroJSON{
				Titulo:     l.GetTitulo(),
				Autor:      l.GetAutor(),
				Categoria:  nombreCat,
				Disponible: disponible,
			})
		}
	} else {
		catID, parseErr := strconv.ParseUint(categoriaIDStr, 10, 32)
		if parseErr != nil {
			http.Error(w, `{"error":"categoria_id inválido"}`, http.StatusBadRequest)
			return
		}
		libros, err := h.libroService.BuscarPorCategoria(uint(catID))
		if err != nil {
			log.Printf("Error al buscar libros por categoría: %v", err)
			http.Error(w, `{"error":"Error al buscar libros"}`, http.StatusInternalServerError)
			return
		}
		nombreCat := "Sin categoría"
		if nombre, ok := catMap[uint(catID)]; ok {
			nombreCat = nombre
		}
		for _, l := range libros {
			disponible := "SI"
			if !l.GetDisponible() {
				disponible = "NO"
			}
			librosJSON = append(librosJSON, LibroJSON{
				Titulo:     l.GetTitulo(),
				Autor:      l.GetAutor(),
				Categoria:  nombreCat,
				Disponible: disponible,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(librosJSON); err != nil {
		log.Printf("Error al codificar JSON: %v", err)
		http.Error(w, `{"error":"Error al generar respuesta"}`, http.StatusInternalServerError)
	}
}

// ListarPrestamosPorUsuario devuelve los préstamos de un usuario recibido en JSON.
// Endpoint: POST /api/prestamos/por-usuario
// Body: {"nombre_usuario":"Wendy Soto"}
func (h *APIHandler) ListarPrestamosPorUsuario(w http.ResponseWriter, r *http.Request) {
	var req PrestamoPorUsuarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON inválido"}`, http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(req.NombreUsuario)
	if nombre == "" {
		http.Error(w, `{"error":"nombre_usuario es obligatorio"}`, http.StatusBadRequest)
		return
	}

	usuarios, err := h.usuarioService.ListarUsuarios()
	if err != nil {
		log.Printf("Error al listar usuarios: %v", err)
		http.Error(w, `{"error":"Error al obtener usuarios"}`, http.StatusInternalServerError)
		return
	}

	var usuarioEncontrado *models.Usuario
	for _, usr := range usuarios {
		if strings.EqualFold(strings.TrimSpace(usr.GetNombre()), nombre) {
			usuarioEncontrado = usr
			break
		}
	}
	if usuarioEncontrado == nil {
		http.Error(w, `{"error":"usuario no encontrado"}`, http.StatusNotFound)
		return
	}

	prestamos, err := h.prestamoService.ObtenerPorUsuario(usuarioEncontrado.GetID())
	if err != nil {
		log.Printf("Error al obtener préstamos del usuario: %v", err)
		http.Error(w, `{"error":"Error al obtener préstamos"}`, http.StatusInternalServerError)
		return
	}

	var prestamosJSON []PrestamoJSON
	for _, p := range prestamos {
		libroTitulo := "Desconocido"
		if libro, err := h.libroService.ObtenerLibro(p.GetLibroID()); err == nil {
			libroTitulo = libro.GetTitulo()
		}

		prestamosJSON = append(prestamosJSON, PrestamoJSON{
			ID:              p.GetID(),
			Libro:           libroTitulo,
			Usuario:         usuarioEncontrado.GetNombre(),
			FechaPrestamo:   p.GetFechaPrestamo().Format("2006-01-02"),
			FechaDevolucion: p.FechaDevolucionStr(),
			Estado:          string(p.GetEstado()),
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(prestamosJSON); err != nil {
		log.Printf("Error al codificar JSON: %v", err)
		http.Error(w, `{"error":"Error al generar respuesta"}`, http.StatusInternalServerError)
	}
}
