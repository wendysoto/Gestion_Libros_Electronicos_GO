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
// INTERFAZ: CategoriaRepository
// Define las operaciones CRUD para la entidad Categoria.
// ============================================================
type CategoriaRepository interface {
	Crear(categoria *models.Categoria) error
	ObtenerPorID(id uint) (*models.Categoria, error)
	ObtenerTodas() ([]*models.Categoria, error)
	Actualizar(categoria *models.Categoria) error
	Eliminar(id uint) error
}

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
	BuscarPorCategoria(categoriaID uint) ([]*models.Libro, error)
}

// ============================================================
// INTERFAZ: UsuarioRepository
// Define las operaciones CRUD para la entidad Usuario.
// ============================================================
type UsuarioRepository interface {
	Crear(usuario *models.Usuario) error
	ObtenerPorID(id uint) (*models.Usuario, error)
	ObtenerTodos() ([]*models.Usuario, error)
	Actualizar(usuario *models.Usuario) error
	Eliminar(id uint) error
}

// ============================================================
// INTERFAZ: PrestamoRepository
// Define las operaciones para gestionar préstamos de libros.
// ============================================================
type PrestamoRepository interface {
	Crear(prestamo *models.Prestamo) error
	ObtenerPorID(id uint) (*models.Prestamo, error)
	ObtenerTodos() ([]*models.Prestamo, error)
	Actualizar(prestamo *models.Prestamo) error
	ObtenerActivosPorUsuario(usuarioID uint) ([]*models.Prestamo, error)
	ObtenerPorUsuario(usuarioID uint) ([]*models.Prestamo, error)
}

// ============================================================
// INTERFAZ: Entidad
// Interfaz genérica que demuestra polimorfismo: cualquier entidad
// del sistema que tenga un ID y representación en texto la cumple.
// Tanto Libro, Usuario, Categoria y Prestamo satisfacen esta interfaz.
// ============================================================
type Entidad interface {
	GetID() uint
	String() string
}
