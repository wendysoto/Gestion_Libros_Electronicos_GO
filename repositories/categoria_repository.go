// Package repositories - Implementación concreta de CategoriaRepository para PostgreSQL.
// POLIMORFISMO: CategoriaRepositoryPostgres satisface la interfaz CategoriaRepository.
package repositories

import (
	"fmt"
	"time"

	"gestion_libros/db"
	"gestion_libros/models"

	"gorm.io/gorm"
)

// CategoriaRepositoryPostgres implementa interfaces.CategoriaRepository.
type CategoriaRepositoryPostgres struct {
	db *gorm.DB
}

// NuevoCategoriaRepository crea una nueva instancia del repositorio.
func NuevoCategoriaRepository(database *gorm.DB) *CategoriaRepositoryPostgres {
	return &CategoriaRepositoryPostgres{db: database}
}

func (r *CategoriaRepositoryPostgres) Crear(cat *models.Categoria) error {
	catDB := &db.CategoriaDB{
		Nombre:      cat.GetNombre(),
		Descripcion: cat.GetDescripcion(),
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}
	resultado := r.db.Create(catDB)
	if resultado.Error != nil {
		return models.NuevoAppError("CategoriaRepository.Crear", "categoria", resultado.Error)
	}
	cat.SetID(catDB.ID)
	return nil
}

func (r *CategoriaRepositoryPostgres) ObtenerPorID(id uint) (*models.Categoria, error) {
	var catDB db.CategoriaDB
	resultado := r.db.First(&catDB, id)
	if resultado.Error != nil {
		if resultado.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: categoría con ID %d", models.ErrNoEncontrado, id)
		}
		return nil, models.NuevoAppError("CategoriaRepository.ObtenerPorID", "categoria", resultado.Error)
	}
	return convertirACategoriaModelo(&catDB), nil
}

func (r *CategoriaRepositoryPostgres) ObtenerTodas() ([]*models.Categoria, error) {
	var categoriasDB []db.CategoriaDB
	resultado := r.db.Order("nombre ASC").Find(&categoriasDB)
	if resultado.Error != nil {
		return nil, models.NuevoAppError("CategoriaRepository.ObtenerTodas", "categoria", resultado.Error)
	}
	categorias := make([]*models.Categoria, len(categoriasDB))
	for i, catDB := range categoriasDB {
		categorias[i] = convertirACategoriaModelo(&catDB)
	}
	return categorias, nil
}

func (r *CategoriaRepositoryPostgres) Actualizar(cat *models.Categoria) error {
	resultado := r.db.Model(&db.CategoriaDB{}).Where("id = ?", cat.GetID()).Updates(map[string]interface{}{
		"nombre":      cat.GetNombre(),
		"descripcion": cat.GetDescripcion(),
		"updated_at":  time.Now().Format("2006-01-02 15:04:05"),
	})
	if resultado.Error != nil {
		return models.NuevoAppError("CategoriaRepository.Actualizar", "categoria", resultado.Error)
	}
	return nil
}

func (r *CategoriaRepositoryPostgres) Eliminar(id uint) error {
	resultado := r.db.Delete(&db.CategoriaDB{}, id)
	if resultado.Error != nil {
		return models.NuevoAppError("CategoriaRepository.Eliminar", "categoria", resultado.Error)
	}
	if resultado.RowsAffected == 0 {
		return fmt.Errorf("%w: categoría con ID %d", models.ErrNoEncontrado, id)
	}
	return nil
}

// convertirACategoriaModelo transforma un registro de BD a modelo de dominio.
func convertirACategoriaModelo(catDB *db.CategoriaDB) *models.Categoria {
	cat := &models.Categoria{}
	cat.SetID(catDB.ID)
	_ = cat.SetNombre(catDB.Nombre)
	cat.SetDescripcion(catDB.Descripcion)
	if catDB.CreatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", catDB.CreatedAt); err == nil {
			cat.SetCreatedAt(t)
		}
	}
	if catDB.UpdatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", catDB.UpdatedAt); err == nil {
			cat.SetUpdatedAt(t)
		}
	}
	return cat
}
