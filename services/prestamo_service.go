// Package services - Servicio de lógica de negocio para Préstamos.
// Reglas de negocio:
// - Control del límite de préstamos activos por usuario (máximo 3).
// - Actualización del contador de préstamos del usuario.
package services

import (
	"fmt"

	"gestion_libros/interfaces"
	"gestion_libros/models"
)

// PrestamoService orquesta la lógica de préstamo y devolución.
type PrestamoService struct {
	prestamoRepo interfaces.PrestamoRepository
	libroRepo    interfaces.LibroRepository
	usuarioRepo  interfaces.UsuarioRepository
}

// NuevoPrestamoService crea el servicio inyectando las tres dependencias.
func NuevoPrestamoService(
	prestamoRepo interfaces.PrestamoRepository,
	libroRepo interfaces.LibroRepository,
	usuarioRepo interfaces.UsuarioRepository,
) *PrestamoService {
	return &PrestamoService{
		prestamoRepo: prestamoRepo,
		libroRepo:    libroRepo,
		usuarioRepo:  usuarioRepo,
	}
}

// CrearPrestamo registra un nuevo préstamo validando las reglas de negocio:
// 1. El libro debe existir.
// 2. El usuario debe existir y no haber alcanzado el límite de 3 préstamos.
func (s *PrestamoService) CrearPrestamo(libroID, usuarioID uint) (*models.Prestamo, error) {
	// Validar que el libro existe
	_, err := s.libroRepo.ObtenerPorID(libroID)
	if err != nil {
		return nil, fmt.Errorf("libro no encontrado: %w", err)
	}

	// Validar que el usuario existe y puede pedir préstamos
	usuario, err := s.usuarioRepo.ObtenerPorID(usuarioID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}
	if !usuario.PuedePrestar() {
		return nil, fmt.Errorf("%w: el usuario '%s' ya tiene %d préstamos activos (máximo %d)",
			models.ErrOperacionInvalida, usuario.GetNombre(),
			usuario.GetPrestados(), models.MaxPrestamosActivos)
	}

	// Crear el préstamo
	prestamo, err := models.NuevoPrestamo(libroID, usuarioID)
	if err != nil {
		return nil, err
	}
	if err := s.prestamoRepo.Crear(prestamo); err != nil {
		return nil, err
	}

	// Incrementar el contador de préstamos del usuario
	if err := usuario.SetPrestados(usuario.GetPrestados() + 1); err != nil {
		return nil, err
	}
	if err := s.usuarioRepo.Actualizar(usuario); err != nil {
		return nil, fmt.Errorf("error al actualizar préstamos del usuario: %w", err)
	}

	return prestamo, nil
}

// DevolverPrestamo registra la devolución de un libro:
// 1. Marca el préstamo como devuelto con la fecha actual.
// 2. Decrementa el contador de préstamos del usuario.
func (s *PrestamoService) DevolverPrestamo(prestamoID uint) error {
	// Obtener el préstamo
	prestamo, err := s.prestamoRepo.ObtenerPorID(prestamoID)
	if err != nil {
		return fmt.Errorf("préstamo no encontrado: %w", err)
	}

	// Registrar la devolución (valida que no esté ya devuelto)
	if err := prestamo.RegistrarDevolucion(); err != nil {
		return err
	}
	if err := s.prestamoRepo.Actualizar(prestamo); err != nil {
		return fmt.Errorf("error al actualizar préstamo: %w", err)
	}

	// Decrementar el contador de préstamos del usuario
	usuario, err := s.usuarioRepo.ObtenerPorID(prestamo.GetUsuarioID())
	if err != nil {
		return fmt.Errorf("error al obtener usuario del préstamo: %w", err)
	}
	nuevaCantidad := usuario.GetPrestados() - 1
	if nuevaCantidad < 0 {
		nuevaCantidad = 0
	}
	if err := usuario.SetPrestados(nuevaCantidad); err != nil {
		return err
	}
	if err := s.usuarioRepo.Actualizar(usuario); err != nil {
		return fmt.Errorf("error al actualizar préstamos del usuario: %w", err)
	}

	return nil
}

// ListarPrestamos retorna todos los préstamos del sistema.
func (s *PrestamoService) ListarPrestamos() ([]*models.Prestamo, error) {
	return s.prestamoRepo.ObtenerTodos()
}

// ObtenerPrestamo busca un préstamo por ID.
func (s *PrestamoService) ObtenerPrestamo(id uint) (*models.Prestamo, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: el ID debe ser mayor a cero", models.ErrValidacion)
	}
	return s.prestamoRepo.ObtenerPorID(id)
}

// ObtenerPorUsuario retorna todos los préstamos de un usuario.
func (s *PrestamoService) ObtenerPorUsuario(usuarioID uint) ([]*models.Prestamo, error) {
	if usuarioID == 0 {
		return nil, fmt.Errorf("%w: el ID debe ser mayor a cero", models.ErrValidacion)
	}
	return s.prestamoRepo.ObtenerPorUsuario(usuarioID)
}
