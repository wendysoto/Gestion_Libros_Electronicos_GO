// Package repositories contiene las implementaciones concretas de las interfaces
// de repositorio. Cada repositorio interactúa con la base de datos PostgreSQL
// a través de GORM.
//
// POLIMORFISMO: LibroRepositoryPostgres implementa la interfaz LibroRepository.
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
// Recibe la conexión a la base de datos como dependencia (inyección de dependencias).
func NuevoLibroRepository(database *gorm.DB) *LibroRepositoryPostgres {
	return &LibroRepositoryPostgres{db: database}
}

// Crear inserta un nuevo libro en la base de datos.
// Convierte el modelo de dominio a modelo de base de datos antes de persistir.
func (r *LibroRepositoryPostgres) Crear(libro *models.Libro) error {
	libroDB := &db.LibroDB{
		Titulo:     libro.GetTitulo(),
		Autor:      libro.GetAutor(),
		Categoria:  libro.GetCategoria(),
		ISBN:       libro.GetISBN(),
		Formato:    string(libro.GetFormato()),
		Disponible: libro.GetDisponible(),
		CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:  time.Now().Format("2006-01-02 15:04:05"),
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
// Retorna ErrNoEncontrado si el libro no existe.
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
	libroDB := &db.LibroDB{
		ID:         libro.GetID(),
		Titulo:     libro.GetTitulo(),
		Autor:      libro.GetAutor(),
		Categoria:  libro.GetCategoria(),
		ISBN:       libro.GetISBN(),
		Formato:    string(libro.GetFormato()),
		Disponible: libro.GetDisponible(),
		UpdatedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	resultado := r.db.Save(libroDB)
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

// BuscarPorTitulo busca libros cuyo título contenga el texto dado.
// Usa ILIKE para búsqueda case-insensitive en PostgreSQL.
func (r *LibroRepositoryPostgres) BuscarPorTitulo(titulo string) ([]*models.Libro, error) {
	var librosDB []db.LibroDB
	resultado := r.db.Where("titulo ILIKE ?", "%"+titulo+"%").Find(&librosDB)
	if resultado.Error != nil {
		return nil, models.NuevoAppError("LibroRepository.BuscarPorTitulo", "libro", resultado.Error)
	}

	libros := make([]*models.Libro, len(librosDB))
	for i, libroDB := range librosDB {
		libros[i] = convertirALibroModelo(&libroDB)
	}
	return libros, nil
}

// BuscarPorAutor busca libros por autor (búsqueda parcial).
func (r *LibroRepositoryPostgres) BuscarPorAutor(autor string) ([]*models.Libro, error) {
	var librosDB []db.LibroDB
	resultado := r.db.Where("autor ILIKE ?", "%"+autor+"%").Find(&librosDB)
	if resultado.Error != nil {
		return nil, models.NuevoAppError("LibroRepository.BuscarPorAutor", "libro", resultado.Error)
	}

	libros := make([]*models.Libro, len(librosDB))
	for i, libroDB := range librosDB {
		libros[i] = convertirALibroModelo(&libroDB)
	}
	return libros, nil
}

// BuscarPorCategoria busca libros por categoría.
func (r *LibroRepositoryPostgres) BuscarPorCategoria(categoria string) ([]*models.Libro, error) {
	var librosDB []db.LibroDB
	resultado := r.db.Where("categoria ILIKE ?", "%"+categoria+"%").Find(&librosDB)
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
