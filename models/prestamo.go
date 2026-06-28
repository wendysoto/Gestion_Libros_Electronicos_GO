// Package models - Modelo Prestamo
// Demuestra ENCAPSULACION y MANEJO DE ERRORES en la gestión de préstamos.
package models

import (
	"fmt"
	"time"
)

// EstadoPrestamo define los estados posibles de un préstamo.
type EstadoPrestamo string

const (
	EstadoActivo   EstadoPrestamo = "activo"
	EstadoDevuelto EstadoPrestamo = "devuelto"
)

// Prestamo representa un préstamo de libro a un usuario.
// Campos encapsulados con validación en los setters.
type Prestamo struct {
	id              uint
	libroID         uint
	usuarioID       uint
	fechaPrestamo   time.Time
	fechaDevolucion *time.Time
	estado          EstadoPrestamo
	createdAt       time.Time
	updatedAt       time.Time
}

// --- Métodos Getter ---

func (p *Prestamo) GetID() uint                  { return p.id }
func (p *Prestamo) GetLibroID() uint             { return p.libroID }
func (p *Prestamo) GetUsuarioID() uint           { return p.usuarioID }
func (p *Prestamo) GetFechaPrestamo() time.Time  { return p.fechaPrestamo }
func (p *Prestamo) GetFechaDevolucion() *time.Time { return p.fechaDevolucion }
func (p *Prestamo) GetEstado() EstadoPrestamo    { return p.estado }
func (p *Prestamo) GetCreatedAt() time.Time      { return p.createdAt }
func (p *Prestamo) GetUpdatedAt() time.Time      { return p.updatedAt }

// FechaDevolucionStr retorna la fecha de devolución formateada o "Pendiente".
func (p *Prestamo) FechaDevolucionStr() string {
	if p.fechaDevolucion == nil {
		return "Pendiente"
	}
	return p.fechaDevolucion.Format("2006-01-02")
}

// --- Métodos Setter ---

func (p *Prestamo) SetID(id uint)            { p.id = id }
func (p *Prestamo) SetCreatedAt(t time.Time) { p.createdAt = t }
func (p *Prestamo) SetUpdatedAt(t time.Time) { p.updatedAt = t }

// SetLibroID establece el ID del libro prestado.
func (p *Prestamo) SetLibroID(id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: el ID del libro debe ser mayor a cero", ErrValidacion)
	}
	p.libroID = id
	return nil
}

// SetUsuarioID establece el ID del usuario que recibe el préstamo.
func (p *Prestamo) SetUsuarioID(id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: el ID del usuario debe ser mayor a cero", ErrValidacion)
	}
	p.usuarioID = id
	return nil
}

// SetEstado cambia el estado del préstamo.
func (p *Prestamo) SetEstado(estado EstadoPrestamo) error {
	switch estado {
	case EstadoActivo, EstadoDevuelto:
		p.estado = estado
		p.updatedAt = time.Now()
		return nil
	default:
		return fmt.Errorf("%w: estado '%s' no válido", ErrValidacion, estado)
	}
}

// SetFechaPrestamo establece la fecha de préstamo.
func (p *Prestamo) SetFechaPrestamo(t time.Time) {
	p.fechaPrestamo = t
	p.updatedAt = time.Now()
}

// SetFechaDevolucion establece la fecha de devolución.
func (p *Prestamo) SetFechaDevolucion(t *time.Time) {
	p.fechaDevolucion = t
	p.updatedAt = time.Now()
}

// RegistrarDevolucion marca el préstamo como devuelto y registra la fecha.
// Retorna error si el préstamo ya fue devuelto previamente.
func (p *Prestamo) RegistrarDevolucion() error {
	if p.estado == EstadoDevuelto {
		return fmt.Errorf("%w: el préstamo ya fue devuelto", ErrOperacionInvalida)
	}
	ahora := time.Now()
	p.fechaDevolucion = &ahora
	p.estado = EstadoDevuelto
	p.updatedAt = ahora
	return nil
}

// NuevoPrestamo crea un préstamo con estado activo y fecha actual.
func NuevoPrestamo(libroID, usuarioID uint) (*Prestamo, error) {
	p := &Prestamo{
		fechaPrestamo: time.Now(),
		estado:        EstadoActivo,
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
	}
	if err := p.SetLibroID(libroID); err != nil {
		return nil, err
	}
	if err := p.SetUsuarioID(usuarioID); err != nil {
		return nil, err
	}
	return p, nil
}

// String implementa fmt.Stringer.
func (p *Prestamo) String() string {
	return fmt.Sprintf("[%d] Libro:%d → Usuario:%d (%s)", p.id, p.libroID, p.usuarioID, p.estado)
}
