package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// TipoUsuario representa los roles posibles de un usuario en el sistema.
type TipoUsuario string

const (
	TipoAdmin  TipoUsuario = "admin"
	TipoLector TipoUsuario = "lector"
)

// Usuario representa un usuario registrado en el sistema de biblioteca.
// Campos encapsulados con acceso controlado mediante getters y setters.
type Usuario struct {
	id        uint
	nombre    string
	email     string
	tipo      TipoUsuario
	createdAt time.Time
	updatedAt time.Time
}

// --- Métodos Getter ---

func (u *Usuario) GetID() uint            { return u.id }
func (u *Usuario) GetNombre() string       { return u.nombre }
func (u *Usuario) GetEmail() string        { return u.email }
func (u *Usuario) GetTipo() TipoUsuario    { return u.tipo }
func (u *Usuario) GetCreatedAt() time.Time { return u.createdAt }
func (u *Usuario) GetUpdatedAt() time.Time { return u.updatedAt }

// --- Métodos Setter con validación ---

// SetNombre valida y establece el nombre del usuario.
// Retorna error si el nombre está vacío o excede 255 caracteres.
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

// SetEmail valida el formato del email y lo establece.
// Utiliza una expresión regular para verificar el formato correcto.
func (u *Usuario) SetEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return fmt.Errorf("%w: el email no puede estar vacío", ErrValidacion)
	}

	// Validación de formato de email con expresión regular
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("%w: formato de email inválido '%s'", ErrValidacion, email)
	}

	u.email = email
	u.updatedAt = time.Now()
	return nil
}

// SetTipo establece el rol del usuario en el sistema.
// Solo acepta valores válidos: 'admin' o 'lector'.
func (u *Usuario) SetTipo(tipo TipoUsuario) error {
	switch tipo {
	case TipoAdmin, TipoLector:
		u.tipo = tipo
		u.updatedAt = time.Now()
		return nil
	default:
		return fmt.Errorf("%w: tipo de usuario '%s' no válido (use 'admin' o 'lector')", ErrValidacion, tipo)
	}
}

// SetID establece el ID del usuario (uso interno del repositorio).
func (u *Usuario) SetID(id uint) {
	u.id = id
}

// SetCreatedAt establece la fecha de creación (uso interno).
func (u *Usuario) SetCreatedAt(t time.Time) {
	u.createdAt = t
}

// SetUpdatedAt establece la fecha de actualización (uso interno).
func (u *Usuario) SetUpdatedAt(t time.Time) {
	u.updatedAt = t
}

// NuevoUsuario crea un usuario validando todos los campos requeridos.
// Demuestra encapsulación: obliga a usar el constructor para crear instancias válidas.
func NuevoUsuario(nombre, email string, tipo TipoUsuario) (*Usuario, error) {
	usuario := &Usuario{
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	if err := usuario.SetNombre(nombre); err != nil {
		return nil, err
	}
	if err := usuario.SetEmail(email); err != nil {
		return nil, err
	}
	if err := usuario.SetTipo(tipo); err != nil {
		return nil, err
	}

	return usuario, nil
}

// String implementa fmt.Stringer para representación legible del usuario.
func (u *Usuario) String() string {
	return fmt.Sprintf("[%d] %s (%s) - %s", u.id, u.nombre, u.email, u.tipo)
}
