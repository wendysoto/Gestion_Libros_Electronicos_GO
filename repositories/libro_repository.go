// Package repositories - Implementación concreta de LibroRepository para PostgreSQL.
// POLIMORFISMO: LibroRepositoryPostgres satisface la interfaz LibroRepository.
// Si en el futuro se necesita cambiar la base de datos (ej: MySQL, MongoDB),
// solo se necesita crear una nueva implementación que satisfaga la misma interfaz.
package repositories

import (
	"fmt"
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
func NuevoLibroRepository(database *gorm.DB) *LibroRepositoryPostgres {
	return &LibroRepositoryPostgres{db: database}
}

// Crear inserta un nuevo libro en la base de datos.
func (r *LibroRepositoryPostgres) Crear(libro *models.Libro) error {
	libroDB := &db.LibroDB{
		Titulo:      libro.GetTitulo(),
		Autor:       libro.GetAutor(),
		CategoriaID: libro.GetCategoriaID(),
		ISBN:        libro.GetISBN(),
		Formato:     string(libro.GetFormato()),
		Disponible:  libro.GetDisponible(),
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}
	resultado := r.db.Create(libroDB)
	if resultado.Error != nil {
		return models.NuevoAppError("LibroRepository.Crear", "libro",
			fmt.Errorf("%w: %v", models.ErrConexionDB, resultado.Error))
	}
	libro.SetID(libroDB.ID)
	return nil
}

// ObtenerPorID busca un libro por su ID en la base de datos.
func (r *LibroRepositoryPostgres) ObtenerPorID(id uint) (*models.Libro, error) {
	var libroDB db.LibroDB
	resultado := r.db.First(&libroDB, id)
	if resultado.Error != nil {
		if resultado.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: libro con ID %d", models.ErrNoEncontrado, id)
		}
		return nil, models.NuevoAppError("LibroRepository.ObtenerPorID", "libro", resultado.Error)
	}
	return convertirALibroModelo(&libroDB), nil
}

// ObtenerTodos retorna todos los libros registrados en la base de datos.
func (r *LibroRepositoryPostgres) ObtenerTodos() ([]*models.Libro, error) {
	var librosDB []db.LibroDB
	resultado := r.db.Find(&librosDB)
	if resultado.Error != nil {
		return nil, models.NuevoAppError("LibroRepository.ObtenerTodos", "libro", resultado.Error)
	}
	libros := make([]*models.Libro, len(librosDB))
	for i, libroDB := range librosDB {
		libros[i] = convertirALibroModelo(&libroDB)
	}
	return libros, nil
}

// Actualizar modifica un libro existente en la base de datos.
func (r *LibroRepositoryPostgres) Actualizar(libro *models.Libro) error {
	resultado := r.db.Model(&db.LibroDB{}).Where("id = ?", libro.GetID()).Updates(map[string]interface{}{
		"titulo":       libro.GetTitulo(),
		"autor":        libro.GetAutor(),
		"categoria_id": libro.GetCategoriaID(),
		"isbn":         libro.GetISBN(),
		"formato":      string(libro.GetFormato()),
		"disponible":   libro.GetDisponible(),
		"updated_at":   time.Now().Format("2006-01-02 15:04:05"),
	})
	if resultado.Error != nil {
		return models.NuevoAppError("LibroRepository.Actualizar", "libro", resultado.Error)
	}
	return nil
}

// Eliminar borra un libro de la base de datos por su ID.
func (r *LibroRepositoryPostgres) Eliminar(id uint) error {
	resultado := r.db.Delete(&db.LibroDB{}, id)
	if resultado.Error != nil {
		return models.NuevoAppError("LibroRepository.Eliminar", "libro", resultado.Error)
	}
	if resultado.RowsAffected == 0 {
		return fmt.Errorf("%w: libro con ID %d", models.ErrNoEncontrado, id)
	}
	return nil
}

// BuscarPorCategoria retorna los libros de una categoría específica.
func (r *LibroRepositoryPostgres) BuscarPorCategoria(categoriaID uint) ([]*models.Libro, error) {
	var librosDB []db.LibroDB
	resultado := r.db.Where("categoria_id = ?", categoriaID).Find(&librosDB)
	if resultado.Error != nil {
		return nil, models.NuevoAppError("LibroRepository.BuscarPorCategoria", "libro", resultado.Error)
	}
	libros := make([]*models.Libro, len(librosDB))
	for i, libroDB := range librosDB {
		libros[i] = convertirALibroModelo(&libroDB)
	}
	return libros, nil
}

// convertirALibroModelo transforma un registro de BD a modelo de dominio.
// Mantiene la separación entre la capa de datos y la capa de dominio.
func convertirALibroModelo(libroDB *db.LibroDB) *models.Libro {
	libro := &models.Libro{}
	libro.SetID(libroDB.ID)
	_ = libro.SetTitulo(libroDB.Titulo)
	_ = libro.SetAutor(libroDB.Autor)
	_ = libro.SetCategoriaID(libroDB.CategoriaID)
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
