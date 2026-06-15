// Package services contiene la lógica de negocio del sistema.
// Los servicios actúan como intermediarios entre los handlers (HTTP)
// y los repositorios (base de datos), aplicando reglas de negocio.
//
// Aquí se demuestra el uso de INTERFACES como parámetros:
// los servicios reciben interfaces de repositorio, no implementaciones
// concretas. Esto permite desacoplar la lógica de negocio del almacenamiento.
package services

import (
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

// ListarLibros retorna todos los libros del catálogo desde la base de datos.
// Delega la consulta al repositorio a través de la interfaz.
func (s *LibroService) ListarLibros() ([]*models.Libro, error) {
	return s.repo.ObtenerTodos()
}
