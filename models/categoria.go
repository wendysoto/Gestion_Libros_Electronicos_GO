// Package models - Modelo Categoria
// Demuestra ENCAPSULACION: campos privados accesibles solo mediante getters/setters.
package models

import (
	"fmt"
	"strings"
	"time"
)

// Categoria representa una categoría de libros en el sistema.
// Los campos son no exportados (minúscula) para garantizar encapsulación;
// el acceso se realiza exclusivamente a través de métodos getter/setter.
type Categoria struct {
	id          uint
	nombre      string
	descripcion string
	createdAt   time.Time
	updatedAt   time.Time
}

// --- Métodos Getter ---

func (c *Categoria) GetID() uint            { return c.id }
func (c *Categoria) GetNombre() string      { return c.nombre }
func (c *Categoria) GetDescripcion() string { return c.descripcion }
func (c *Categoria) GetCreatedAt() time.Time { return c.createdAt }
func (c *Categoria) GetUpdatedAt() time.Time { return c.updatedAt }

// --- Métodos Setter (validan antes de asignar) ---

// SetNombre establece el nombre de la categoría.
// Retorna error si el nombre está vacío o excede 100 caracteres.
func (c *Categoria) SetNombre(nombre string) error {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return fmt.Errorf("%w: el nombre de la categoría no puede estar vacío", ErrValidacion)
	}
	if len(nombre) > 100 {
		return fmt.Errorf("%w: el nombre no puede exceder 100 caracteres", ErrValidacion)
	}
	c.nombre = nombre
	c.updatedAt = time.Now()
	return nil
}

// SetDescripcion establece la descripción de la categoría.
func (c *Categoria) SetDescripcion(descripcion string) {
	c.descripcion = strings.TrimSpace(descripcion)
	c.updatedAt = time.Now()
}

// SetID establece el ID (uso interno por el repositorio).
func (c *Categoria) SetID(id uint) { c.id = id }

// SetCreatedAt establece la fecha de creación (uso interno).
func (c *Categoria) SetCreatedAt(t time.Time) { c.createdAt = t }

// SetUpdatedAt establece la fecha de actualización (uso interno).
func (c *Categoria) SetUpdatedAt(t time.Time) { c.updatedAt = t }

// NuevaCategoria es un constructor que crea una categoría validada.
// Obliga a pasar por los setters, demostrando encapsulación.
func NuevaCategoria(nombre, descripcion string) (*Categoria, error) {
	cat := &Categoria{
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
	if err := cat.SetNombre(nombre); err != nil {
		return nil, err
	}
	cat.SetDescripcion(descripcion)
	return cat, nil
}

// String implementa la interfaz fmt.Stringer.
func (c *Categoria) String() string {
	return fmt.Sprintf("[%d] %s", c.id, c.nombre)
}
