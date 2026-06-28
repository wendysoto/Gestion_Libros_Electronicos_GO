// Package models - Modelo Usuario
// Demuestra ENCAPSULACION: campos privados con getters/setters que validan datos.
package models

import (
	"fmt"
	"strings"
	"time"
)

// TipoUsuario define los tipos válidos de usuario del sistema.
type TipoUsuario string

const (
	TipoAdmin  TipoUsuario = "admin"
	TipoLector TipoUsuario = "lector"
)

// MaxPrestamosActivos es el límite de libros que un usuario puede tener
// prestados simultáneamente. Se utiliza como regla de negocio.
const MaxPrestamosActivos = 3

// Usuario representa a un usuario registrado en el sistema.
// Campos encapsulados; incluye 'prestados' para rastrear la cantidad
// de préstamos activos del usuario.
type Usuario struct {
	id        uint
	nombre    string
	email     string
	tipo      TipoUsuario
	prestados int
	createdAt time.Time
	updatedAt time.Time
}

// --- Métodos Getter ---

func (u *Usuario) GetID() uint            { return u.id }
func (u *Usuario) GetNombre() string      { return u.nombre }
func (u *Usuario) GetEmail() string       { return u.email }
func (u *Usuario) GetTipo() TipoUsuario   { return u.tipo }
func (u *Usuario) GetPrestados() int      { return u.prestados }
func (u *Usuario) GetCreatedAt() time.Time { return u.createdAt }
func (u *Usuario) GetUpdatedAt() time.Time { return u.updatedAt }

// --- Métodos Setter ---

// SetNombre valida y establece el nombre del usuario.
func (u *Usuario) SetNombre(nombre string) error {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return fmt.Errorf("%w: el nombre no puede estar vacío", ErrValidacion)
	}
	if len(nombre) > 255 {
		return fmt.Errorf("%w: el nombre no puede exceder 255 caracteres", ErrValidacion)
	}
	u.nombre = nombre
	u.updatedAt = time.Now()
	return nil
}

// SetEmail valida y establece el correo electrónico.
// Verifica que contenga '@' como validación básica de formato.
func (u *Usuario) SetEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("%w: el email no puede estar vacío", ErrValidacion)
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("%w: el email debe contener '@'", ErrValidacion)
	}
	u.email = email
	u.updatedAt = time.Now()
	return nil
}

// SetTipo establece el tipo de usuario.
// Solo acepta 'admin' o 'lector'.
func (u *Usuario) SetTipo(tipo TipoUsuario) error {
	switch tipo {
	case TipoAdmin, TipoLector:
		u.tipo = tipo
		u.updatedAt = time.Now()
		return nil
	default:
		return fmt.Errorf("%w: tipo '%s' no válido (use admin o lector)", ErrValidacion, tipo)
	}
}

// SetPrestados establece la cantidad de préstamos activos.
// No permite valores negativos ni superiores al máximo permitido.
func (u *Usuario) SetPrestados(prestados int) error {
	if prestados < 0 {
		return fmt.Errorf("%w: la cantidad de préstamos no puede ser negativa", ErrValidacion)
	}
	if prestados > MaxPrestamosActivos {
		return fmt.Errorf("%w: no se puede exceder el límite de %d préstamos activos", ErrValidacion, MaxPrestamosActivos)
	}
	u.prestados = prestados
	u.updatedAt = time.Now()
	return nil
}

// PuedePrestar verifica si el usuario tiene capacidad para un nuevo préstamo.
// Regla de negocio: máximo 3 préstamos activos simultáneos.
func (u *Usuario) PuedePrestar() bool {
	return u.prestados < MaxPrestamosActivos
}

func (u *Usuario) SetID(id uint)            { u.id = id }
func (u *Usuario) SetCreatedAt(t time.Time) { u.createdAt = t }
func (u *Usuario) SetUpdatedAt(t time.Time) { u.updatedAt = t }

// NuevoUsuario crea un usuario validado mediante constructores y setters.
func NuevoUsuario(nombre, email string, tipo TipoUsuario) (*Usuario, error) {
	usr := &Usuario{
		prestados: 0,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
	if err := usr.SetNombre(nombre); err != nil {
		return nil, err
	}
	if err := usr.SetEmail(email); err != nil {
		return nil, err
	}
	if err := usr.SetTipo(tipo); err != nil {
		return nil, err
	}
	return usr, nil
}

// String implementa fmt.Stringer.
func (u *Usuario) String() string {
	return fmt.Sprintf("[%d] %s (%s) - Préstamos: %d/%d", u.id, u.nombre, u.tipo, u.prestados, MaxPrestamosActivos)
}
