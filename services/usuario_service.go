package services

import (
	"gestion_libros/interfaces"
	"gestion_libros/models"
)

// UsuarioService maneja la lógica de negocio para usuarios.
// Usa interfaz UsuarioRepository para polimorfismo.
type UsuarioService struct {
	repo interfaces.UsuarioRepository
}

// NuevoUsuarioService crea un servicio con inyección de dependencias.
func NuevoUsuarioService(repo interfaces.UsuarioRepository) *UsuarioService {
	return &UsuarioService{repo: repo}
}

// ListarUsuarios retorna todos los usuarios del sistema desde la base de datos.
func (s *UsuarioService) ListarUsuarios() ([]*models.Usuario, error) {
	return s.repo.ObtenerTodos()
}
