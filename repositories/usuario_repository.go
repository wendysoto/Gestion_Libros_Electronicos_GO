// Package repositories - Implementación concreta de UsuarioRepository para PostgreSQL.
// POLIMORFISMO: UsuarioRepositoryPostgres satisface la interfaz UsuarioRepository.
package repositories

import (
	"fmt"
	"time"

	"gestion_libros/db"
	"gestion_libros/models"

	"gorm.io/gorm"
)

// UsuarioRepositoryPostgres implementa interfaces.UsuarioRepository.
type UsuarioRepositoryPostgres struct {
	db *gorm.DB
}

// NuevoUsuarioRepository crea una instancia del repositorio de usuarios.
func NuevoUsuarioRepository(database *gorm.DB) *UsuarioRepositoryPostgres {
	return &UsuarioRepositoryPostgres{db: database}
}

func (r *UsuarioRepositoryPostgres) Crear(usr *models.Usuario) error {
	usrDB := &db.UsuarioDB{
		Nombre:    usr.GetNombre(),
		Email:     usr.GetEmail(),
		Tipo:      string(usr.GetTipo()),
		Prestados: usr.GetPrestados(),
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	resultado := r.db.Create(usrDB)
	if resultado.Error != nil {
		return models.NuevoAppError("UsuarioRepository.Crear", "usuario", resultado.Error)
	}
	usr.SetID(usrDB.ID)
	return nil
}

func (r *UsuarioRepositoryPostgres) ObtenerPorID(id uint) (*models.Usuario, error) {
	var usrDB db.UsuarioDB
	resultado := r.db.First(&usrDB, id)
	if resultado.Error != nil {
		if resultado.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%w: usuario con ID %d", models.ErrNoEncontrado, id)
		}
		return nil, models.NuevoAppError("UsuarioRepository.ObtenerPorID", "usuario", resultado.Error)
	}
	return convertirAUsuarioModelo(&usrDB), nil
}

func (r *UsuarioRepositoryPostgres) ObtenerTodos() ([]*models.Usuario, error) {
	var usuariosDB []db.UsuarioDB
	resultado := r.db.Find(&usuariosDB)
	if resultado.Error != nil {
		return nil, models.NuevoAppError("UsuarioRepository.ObtenerTodos", "usuario", resultado.Error)
	}
	usuarios := make([]*models.Usuario, len(usuariosDB))
	for i, usrDB := range usuariosDB {
		usuarios[i] = convertirAUsuarioModelo(&usrDB)
	}
	return usuarios, nil
}

func (r *UsuarioRepositoryPostgres) Actualizar(usr *models.Usuario) error {
	resultado := r.db.Model(&db.UsuarioDB{}).Where("id = ?", usr.GetID()).Updates(map[string]interface{}{
		"nombre":     usr.GetNombre(),
		"email":      usr.GetEmail(),
		"tipo":       string(usr.GetTipo()),
		"prestados":  usr.GetPrestados(),
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	})
	if resultado.Error != nil {
		return models.NuevoAppError("UsuarioRepository.Actualizar", "usuario", resultado.Error)
	}
	return nil
}

func (r *UsuarioRepositoryPostgres) Eliminar(id uint) error {
	resultado := r.db.Delete(&db.UsuarioDB{}, id)
	if resultado.Error != nil {
		return models.NuevoAppError("UsuarioRepository.Eliminar", "usuario", resultado.Error)
	}
	if resultado.RowsAffected == 0 {
		return fmt.Errorf("%w: usuario con ID %d", models.ErrNoEncontrado, id)
	}
	return nil
}

// convertirAUsuarioModelo transforma un registro de BD a modelo de dominio.
func convertirAUsuarioModelo(usrDB *db.UsuarioDB) *models.Usuario {
	usr := &models.Usuario{}
	usr.SetID(usrDB.ID)
	_ = usr.SetNombre(usrDB.Nombre)
	_ = usr.SetEmail(usrDB.Email)
	_ = usr.SetTipo(models.TipoUsuario(usrDB.Tipo))
	_ = usr.SetPrestados(usrDB.Prestados)
	if usrDB.CreatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", usrDB.CreatedAt); err == nil {
			usr.SetCreatedAt(t)
		}
	}
	if usrDB.UpdatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", usrDB.UpdatedAt); err == nil {
			usr.SetUpdatedAt(t)
		}
	}
	return usr
}
