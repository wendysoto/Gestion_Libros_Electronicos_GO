// Package repositories contiene las implementaciones concretas de las interfaces
// de repositorio. Cada repositorio interactúa con la base de datos PostgreSQL
// a través de GORM.
//
// POLIMORFISMO: LibroRepositoryPostgres implementa la interfaz LibroRepository.
// Si en el futuro se necesita cambiar la base de datos (ej: MySQL, MongoDB),
// solo se necesita crear una nueva implementación que satisfaga la misma interfaz.
package repositories

import (
	"time"

	"gestion_libros/db"
	"gestion_libros/models"

	"gorm.io/gorm"
)

// LibroRepositoryPostgres implementa la interfaz interfaces.LibroRepository
// usando PostgreSQL como almacenamiento a través de GORM.
type LibroRepositoryPostgres struct {
	db *gorm.DB
}

// NuevoLibroRepository crea una nueva instancia del repositorio de libros.
// Recibe la conexión a la base de datos como dependencia (inyección de dependencias).
func NuevoLibroRepository(database *gorm.DB) *LibroRepositoryPostgres {
	return &LibroRepositoryPostgres{db: database}
}

// ObtenerTodos retorna todos los libros registrados en la base de datos.
// Consulta la tabla 'libros' y convierte cada registro al modelo de dominio.
func (r *LibroRepositoryPostgres) ObtenerTodos() ([]*models.Libro, error) {
	var librosDB []db.LibroDB
	resultado := r.db.Find(&librosDB)
	if resultado.Error != nil {
		return nil, models.NuevoAppError("LibroRepository.ObtenerTodos", "libro", resultado.Error)
	}

	// Convertir cada registro de BD al modelo de dominio encapsulado
	libros := make([]*models.Libro, len(librosDB))
	for i, libroDB := range librosDB {
		libros[i] = convertirALibroModelo(&libroDB)
	}
	return libros, nil
}

// convertirALibroModelo transforma un registro de BD a modelo de dominio.
// Esta función auxiliar mantiene la separación entre la capa de datos
// y la capa de dominio (encapsulación de capas).
func convertirALibroModelo(libroDB *db.LibroDB) *models.Libro {
	libro := &models.Libro{}
	libro.SetID(libroDB.ID)
	_ = libro.SetTitulo(libroDB.Titulo)
	_ = libro.SetAutor(libroDB.Autor)
	_ = libro.SetCategoria(libroDB.Categoria)
	_ = libro.SetISBN(libroDB.ISBN)
	_ = libro.SetFormato(models.FormatoLibro(libroDB.Formato))
	libro.SetDisponible(libroDB.Disponible)

	if libroDB.CreatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", libroDB.CreatedAt); err == nil {
			libro.SetCreatedAt(t)
		}
	}
	if libroDB.UpdatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", libroDB.UpdatedAt); err == nil {
			libro.SetUpdatedAt(t)
		}
	}

	return libro
}
