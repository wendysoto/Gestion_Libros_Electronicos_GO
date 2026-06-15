package repositories

import (
	"time"

	"gestion_libros/db"
	"gestion_libros/models"

	"gorm.io/gorm"
)

// UsuarioRepositoryPostgres implementa la interfaz interfaces.UsuarioRepository.
// Demuestra polimorfismo: esta implementación usa PostgreSQL, pero podría
// existir otra implementación (ej: UsuarioRepositoryMemoria) que satisfaga
// la misma interfaz con almacenamiento en memoria.
type UsuarioRepositoryPostgres struct {
	db *gorm.DB
}

// NuevoUsuarioRepository crea una instancia del repositorio de usuarios.
func NuevoUsuarioRepository(database *gorm.DB) *UsuarioRepositoryPostgres {
	return &UsuarioRepositoryPostgres{db: database}
}

// ObtenerTodos retorna todos los usuarios registrados en la base de datos.
func (r *UsuarioRepositoryPostgres) ObtenerTodos() ([]*models.Usuario, error) {
	var usuariosDB []db.UsuarioDB
	resultado := r.db.Find(&usuariosDB)
	if resultado.Error != nil {
		return nil, models.NuevoAppError("UsuarioRepository.ObtenerTodos", "usuario", resultado.Error)
	}

	usuarios := make([]*models.Usuario, len(usuariosDB))
	for i, usuarioDB := range usuariosDB {
		usuarios[i] = convertirAUsuarioModelo(&usuarioDB)
	}
	return usuarios, nil
}

// convertirAUsuarioModelo transforma un registro de BD a modelo de dominio.
func convertirAUsuarioModelo(usuarioDB *db.UsuarioDB) *models.Usuario {
	usuario := &models.Usuario{}
	usuario.SetID(usuarioDB.ID)
	_ = usuario.SetNombre(usuarioDB.Nombre)
	_ = usuario.SetEmail(usuarioDB.Email)
	_ = usuario.SetTipo(models.TipoUsuario(usuarioDB.Tipo))

	if usuarioDB.CreatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", usuarioDB.CreatedAt); err == nil {
			usuario.SetCreatedAt(t)
		}
	}
	if usuarioDB.UpdatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", usuarioDB.UpdatedAt); err == nil {
			usuario.SetUpdatedAt(t)
		}
	}

	return usuario
}
