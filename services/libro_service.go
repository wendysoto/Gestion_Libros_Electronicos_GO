// Package services contiene la lógica de negocio del sistema.
// Los servicios actúan como intermediarios entre los handlers (HTTP)
// y los repositorios (base de datos), aplicando reglas de negocio.
//
// Aquí se demuestra el uso de INTERFACES como parámetros:
// los servicios reciben interfaces de repositorio, no implementaciones
// concretas. Esto permite desacoplar la lógica de negocio del almacenamiento.
package services

import (
	"fmt"

	"gestion_libros/interfaces"
	"gestion_libros/models"
)

// LibroService maneja la lógica de negocio para libros.
// Recibe una interfaz (no una implementación concreta), lo que permite
// polimorfismo: cualquier implementación de LibroRepository es válida.
type LibroService struct {
	repo interfaces.LibroRepository
}

// NuevoLibroService crea un servicio inyectando la dependencia del repositorio.
// La inyección se hace por interfaz, demostrando desacoplamiento.
func NuevoLibroService(repo interfaces.LibroRepository) *LibroService {
	return &LibroService{repo: repo}
}

// CrearLibro valida los datos y crea un nuevo libro en el sistema.
// Aplica reglas de negocio antes de delegar al repositorio.
func (s *LibroService) CrearLibro(titulo, autor, categoria, isbn string, formato models.FormatoLibro) (*models.Libro, error) {
	// Crear libro con validación mediante constructor (encapsulación)
	libro, err := models.NuevoLibro(titulo, autor, categoria, isbn, formato)
	if err != nil {
		return nil, err
	}

	// Delegar la persistencia al repositorio
	if err := s.repo.Crear(libro); err != nil {
		return nil, err
	}

	return libro, nil
}

// ObtenerLibro busca un libro por ID con manejo de errores adecuado.
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

// ActualizarLibro modifica un libro existente validando los nuevos datos.
func (s *LibroService) ActualizarLibro(id uint, titulo, autor, categoria, isbn string, formato models.FormatoLibro) (*models.Libro, error) {
	libro, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, err
	}

	if titulo != "" {
		if err := libro.SetTitulo(titulo); err != nil {
			return nil, err
		}
	}
	if autor != "" {
		if err := libro.SetAutor(autor); err != nil {
			return nil, err
		}
	}
	if categoria != "" {
		if err := libro.SetCategoria(categoria); err != nil {
			return nil, err
		}
	}
	if isbn != "" {
		if err := libro.SetISBN(isbn); err != nil {
			return nil, err
		}
	}
	if formato != "" {
		if err := libro.SetFormato(formato); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Actualizar(libro); err != nil {
		return nil, err
	}

	return libro, nil
}

// EliminarLibro remueve un libro del catálogo.
func (s *LibroService) EliminarLibro(id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: el ID debe ser mayor a cero", models.ErrValidacion)
	}
	return s.repo.Eliminar(id)
}

// BuscarLibros realiza una búsqueda de libros por el criterio dado.
// Demuestra el uso de la estructura de datos slice y búsqueda flexible.
func (s *LibroService) BuscarLibros(criterio, valor string) ([]*models.Libro, error) {
	switch criterio {
	case "titulo":
		return s.repo.BuscarPorTitulo(valor)
	case "autor":
		return s.repo.BuscarPorAutor(valor)
	case "categoria":
		return s.repo.BuscarPorCategoria(valor)
	default:
		return nil, fmt.Errorf("%w: criterio de búsqueda '%s' no válido", models.ErrValidacion, criterio)
	}
}
