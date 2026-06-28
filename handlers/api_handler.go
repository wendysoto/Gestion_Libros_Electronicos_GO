// Package handlers - API REST en formato JSON.
// Endpoint para listar libros por categoría retornando JSON.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gestion_libros/services"
)

// APIHandler gestiona los endpoints REST que retornan JSON.
type APIHandler struct {
	libroService     *services.LibroService
	categoriaService *services.CategoriaService
}

// NuevoAPIHandler crea una instancia del handler de API.
func NuevoAPIHandler(libroService *services.LibroService, categoriaService *services.CategoriaService) *APIHandler {
	return &APIHandler{
		libroService:     libroService,
		categoriaService: categoriaService,
	}
}

// LibroJSON es la estructura de respuesta JSON para el endpoint de libros.
type LibroJSON struct {
	Titulo     string `json:"titulo"`
	Autor      string `json:"autor"`
	Categoria  string `json:"categoria"`
	Disponible string `json:"disponible"`
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
