// Package interfaces define los contratos (interfaces) que deben cumplir
// las implementaciones del sistema.
//
// POLIMORFISMO EN GO:
// En Go, el polimorfismo se logra mediante interfaces. Una interfaz define
// un conjunto de métodos que un tipo debe implementar. Cualquier tipo que
// implemente todos los métodos de una interfaz satisface automáticamente
// esa interfaz (tipado estructural / duck typing).
//
// Esto permite que diferentes implementaciones (por ejemplo, una que use
// PostgreSQL y otra que use memoria) puedan ser intercambiables, ya que
// ambas satisfacen la misma interfaz.
package interfaces

import "gestion_libros/models"

// ============================================================
// INTERFAZ: LibroRepository
// Define las operaciones CRUD para la entidad Libro.
// Cualquier struct que implemente estos métodos puede actuar
// como repositorio de libros (polimorfismo).
// ============================================================
type LibroRepository interface {
	Crear(libro *models.Libro) error
	ObtenerPorID(id uint) (*models.Libro, error)
	ObtenerTodos() ([]*models.Libro, error)
	Actualizar(libro *models.Libro) error
	Eliminar(id uint) error
	BuscarPorTitulo(titulo string) ([]*models.Libro, error)
	BuscarPorAutor(autor string) ([]*models.Libro, error)
	BuscarPorCategoria(categoria string) ([]*models.Libro, error)
}

// ============================================================
// INTERFAZ: Buscable
// Interfaz genérica que demuestra polimorfismo: cualquier entidad
// que pueda ser buscada por un término implementa esta interfaz.
// ============================================================
type Buscable interface {
	Buscar(termino string) ([]interface{}, error)
}

// ============================================================
// INTERFAZ: Exportable
// Demuestra polimorfismo: permite que diferentes entidades
// se exporten en distintos formatos.
// ============================================================
type Exportable interface {
	ExportarJSON() ([]byte, error)
	ExportarTexto() string
}
