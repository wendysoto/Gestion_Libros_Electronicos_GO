// Package repositories - Implementación concreta de PrestamoRepository para PostgreSQL.
// POLIMORFISMO: PrestamoRepositoryPostgres satisface la interfaz PrestamoRepository.
package repositories

import (
	"fmt"
	"time"

	"gestion_libros/db"
	"gestion_libros/models"

	"gorm.io/gorm"
)

// PrestamoRepositoryPostgres implementa interfaces.PrestamoRepository.
type PrestamoRepositoryPostgres struct {
	db *gorm.DB
}

// NuevoPrestamoRepository crea una instancia del repositorio de préstamos.
func NuevoPrestamoRepository(database *gorm.DB) *PrestamoRepositoryPostgres {
	return &PrestamoRepositoryPostgres{db: database}
}

func (r *PrestamoRepositoryPostgres) Crear(p *models.Prestamo) error {
	pDB := &db.PrestamoDB{
		LibroID:       p.GetLibroID(),
		UsuarioID:     p.GetUsuarioID(),
		FechaPrestamo: p.GetFechaPrestamo(),
		Estado:        string(p.GetEstado()),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	resultado := r.db.Create(pDB)
	if resultado.Error != nil {
		return models.NuevoAppError("PrestamoRepository.Crear", "prestamo", resultado.Error)
	}
	p.SetID(pDB.ID)
	return nil
}

func (r *PrestamoRepositoryPostgres) ObtenerPorID(id uint) (*models.Prestamo, error) {
	var pDB db.PrestamoDB
	resultado := r.db.First(&pDB, id)
	if resultado.Error != nil {
		if resultado.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: préstamo con ID %d", models.ErrNoEncontrado, id)
		}
		return nil, models.NuevoAppError("PrestamoRepository.ObtenerPorID", "prestamo", resultado.Error)
	}
	return convertirAPrestamoModelo(&pDB), nil
}

func (r *PrestamoRepositoryPostgres) ObtenerTodos() ([]*models.Prestamo, error) {
	var prestamosDB []db.PrestamoDB
	resultado := r.db.Order("id DESC").Find(&prestamosDB)
	if resultado.Error != nil {
		return nil, models.NuevoAppError("PrestamoRepository.ObtenerTodos", "prestamo", resultado.Error)
	}
	prestamos := make([]*models.Prestamo, len(prestamosDB))
	for i, pDB := range prestamosDB {
		prestamos[i] = convertirAPrestamoModelo(&pDB)
	}
	return prestamos, nil
}

func (r *PrestamoRepositoryPostgres) Actualizar(p *models.Prestamo) error {
	campos := map[string]interface{}{
		"libro_id":       p.GetLibroID(),
		"usuario_id":     p.GetUsuarioID(),
		"fecha_prestamo": p.GetFechaPrestamo(),
		"estado":         string(p.GetEstado()),
		"updated_at":     time.Now(),
	}
	if p.GetFechaDevolucion() != nil {
		campos["fecha_devolucion"] = p.GetFechaDevolucion()
	}
	resultado := r.db.Model(&db.PrestamoDB{}).Where("id = ?", p.GetID()).Updates(campos)
	if resultado.Error != nil {
		return models.NuevoAppError("PrestamoRepository.Actualizar", "prestamo", resultado.Error)
	}
	return nil
}

// ObtenerActivosPorUsuario retorna los préstamos activos de un usuario específico.
// Se utiliza para validar el límite de préstamos (máximo 3).
func (r *PrestamoRepositoryPostgres) ObtenerActivosPorUsuario(usuarioID uint) ([]*models.Prestamo, error) {
	var prestamosDB []db.PrestamoDB
	resultado := r.db.Where("usuario_id = ? AND estado = ?", usuarioID, "activo").Find(&prestamosDB)
	if resultado.Error != nil {
		return nil, models.NuevoAppError("PrestamoRepository.ObtenerActivosPorUsuario", "prestamo", resultado.Error)
	}
	prestamos := make([]*models.Prestamo, len(prestamosDB))
	for i, pDB := range prestamosDB {
		prestamos[i] = convertirAPrestamoModelo(&pDB)
	}
	return prestamos, nil
}

// ObtenerPorUsuario retorna todos los préstamos de un usuario.
func (r *PrestamoRepositoryPostgres) ObtenerPorUsuario(usuarioID uint) ([]*models.Prestamo, error) {
	var prestamosDB []db.PrestamoDB
	resultado := r.db.Where("usuario_id = ?", usuarioID).Order("id DESC").Find(&prestamosDB)
	if resultado.Error != nil {
		return nil, models.NuevoAppError("PrestamoRepository.ObtenerPorUsuario", "prestamo", resultado.Error)
	}
	prestamos := make([]*models.Prestamo, len(prestamosDB))
	for i, pDB := range prestamosDB {
		prestamos[i] = convertirAPrestamoModelo(&pDB)
	}
	return prestamos, nil
}

// convertirAPrestamoModelo transforma un registro de BD a modelo de dominio.
func convertirAPrestamoModelo(pDB *db.PrestamoDB) *models.Prestamo {
	p := &models.Prestamo{}
	p.SetID(pDB.ID)
	_ = p.SetLibroID(pDB.LibroID)
	_ = p.SetUsuarioID(pDB.UsuarioID)
	_ = p.SetEstado(models.EstadoPrestamo(pDB.Estado))

	p.SetFechaPrestamo(pDB.FechaPrestamo)
	if pDB.FechaDevolucion != nil {
		p.SetFechaDevolucion(pDB.FechaDevolucion)
	}
	p.SetCreatedAt(pDB.CreatedAt)
	p.SetUpdatedAt(pDB.UpdatedAt)
	return p
}
