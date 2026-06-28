// Package services - Servicio de lógica de negocio para Usuarios.
package services

import (
	"fmt"

	"gestion_libros/interfaces"
	"gestion_libros/models"
)

// UsuarioService maneja la lógica de negocio para usuarios.
type UsuarioService struct {
	repo         interfaces.UsuarioRepository
	prestamoRepo interfaces.PrestamoRepository
}

// NuevoUsuarioService crea un servicio inyectando las dependencias.
func NuevoUsuarioService(repo interfaces.UsuarioRepository, prestamoRepo interfaces.PrestamoRepository) *UsuarioService {
	return &UsuarioService{repo: repo, prestamoRepo: prestamoRepo}
}

// CrearUsuario valida y crea un nuevo usuario.
func (s *UsuarioService) CrearUsuario(nombre, email string, tipo models.TipoUsuario) (*models.Usuario, error) {
	usr, err := models.NuevoUsuario(nombre, email, tipo)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Crear(usr); err != nil {
		return nil, err
	}
	return usr, nil
}

// ListarUsuarios retorna todos los usuarios.
func (s *UsuarioService) ListarUsuarios() ([]*models.Usuario, error) {
	return s.repo.ObtenerTodos()
}

// ObtenerUsuario busca un usuario por ID.
func (s *UsuarioService) ObtenerUsuario(id uint) (*models.Usuario, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: el ID debe ser mayor a cero", models.ErrValidacion)
	}
	return s.repo.ObtenerPorID(id)
}

// ActualizarUsuario actualiza los datos de un usuario existente.
func (s *UsuarioService) ActualizarUsuario(usr *models.Usuario) error {
	return s.repo.Actualizar(usr)
}

// EliminarUsuario elimina un usuario solo si no tiene préstamos activos.
func (s *UsuarioService) EliminarUsuario(id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: el ID debe ser mayor a cero", models.ErrValidacion)
	}
	activos, err := s.prestamoRepo.ObtenerActivosPorUsuario(id)
	if err != nil {
		return fmt.Errorf("error al verificar préstamos: %w", err)
	}
	if len(activos) > 0 {
		return fmt.Errorf("no se puede eliminar: el usuario tiene %d préstamo(s) activo(s)", len(activos))
	}
	return s.repo.Eliminar(id)
}
