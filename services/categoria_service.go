// Package services - Servicio de lógica de negocio para Categorías.
// Recibe una interfaz CategoriaRepository, demostrando polimorfismo.
package services

import (
	"fmt"

	"gestion_libros/interfaces"
	"gestion_libros/models"
)

// CategoriaService maneja la lógica de negocio para categorías.
type CategoriaService struct {
	repo interfaces.CategoriaRepository
}

// NuevoCategoriaService crea un servicio inyectando la dependencia.
func NuevoCategoriaService(repo interfaces.CategoriaRepository) *CategoriaService {
	return &CategoriaService{repo: repo}
}

// CrearCategoria valida y crea una nueva categoría.
func (s *CategoriaService) CrearCategoria(nombre, descripcion string) (*models.Categoria, error) {
	cat, err := models.NuevaCategoria(nombre, descripcion)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Crear(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

// ListarCategorias retorna todas las categorías.
func (s *CategoriaService) ListarCategorias() ([]*models.Categoria, error) {
	return s.repo.ObtenerTodas()
}

// ObtenerCategoria busca una categoría por ID.
func (s *CategoriaService) ObtenerCategoria(id uint) (*models.Categoria, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: el ID debe ser mayor a cero", models.ErrValidacion)
	}
	return s.repo.ObtenerPorID(id)
}

// EliminarCategoria elimina una categoría por ID.
func (s *CategoriaService) EliminarCategoria(id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: el ID debe ser mayor a cero", models.ErrValidacion)
	}
	return s.repo.Eliminar(id)
}
