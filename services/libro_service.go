// Package services contiene la lógica de negocio del sistema.
// Los servicios reciben INTERFACES de repositorio (no implementaciones concretas),
// lo que demuestra polimorfismo y desacoplamiento.
package services

import (
	"fmt"

	"gestion_libros/interfaces"
	"gestion_libros/models"
)

// LibroService maneja la lógica de negocio para libros.
// Recibe una interfaz (no una implementación concreta), permitiendo polimorfismo.
type LibroService struct {
	repo interfaces.LibroRepository
}

// NuevoLibroService crea un servicio inyectando la dependencia del repositorio.
func NuevoLibroService(repo interfaces.LibroRepository) *LibroService {
	return &LibroService{repo: repo}
}

// CrearLibro valida los datos y crea un nuevo libro en el sistema.
func (s *LibroService) CrearLibro(titulo, autor string, categoriaID uint, isbn string, formato models.FormatoLibro) (*models.Libro, error) {
	libro, err := models.NuevoLibro(titulo, autor, categoriaID, isbn, formato)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Crear(libro); err != nil {
		return nil, err
	}
	return libro, nil
}

// ObtenerLibro busca un libro por ID con manejo de errores.
func (s *LibroService) ObtenerLibro(id uint) (*models.Libro, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: el ID debe ser mayor a cero", models.ErrValidacion)
	}
	return s.repo.ObtenerPorID(id)
}

// ListarLibros retorna todos los libros del catálogo.
func (s *LibroService) ListarLibros() ([]*models.Libro, error) {
	return s.repo.ObtenerTodos()
}

// EliminarLibro remueve un libro del catálogo.
func (s *LibroService) EliminarLibro(id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: el ID debe ser mayor a cero", models.ErrValidacion)
	}
	return s.repo.Eliminar(id)
}

// ActualizarLibro actualiza los datos de un libro existente.
func (s *LibroService) ActualizarLibro(libro *models.Libro) error {
	return s.repo.Actualizar(libro)
}

// BuscarPorCategoria retorna los libros que pertenecen a una categoría.
func (s *LibroService) BuscarPorCategoria(categoriaID uint) ([]*models.Libro, error) {
	if categoriaID == 0 {
		return nil, fmt.Errorf("%w: el ID de categoría debe ser mayor a cero", models.ErrValidacion)
	}
	return s.repo.BuscarPorCategoria(categoriaID)
}
