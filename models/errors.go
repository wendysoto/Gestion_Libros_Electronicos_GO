// Manejo de errores personalizados del sistema.
// Define errores centinela (sentinel errors) y tipos de error custom
// para clasificar y manejar adecuadamente las excepciones del sistema.
package models

import (
	"errors"
	"fmt"
)

// ============================================================
// ERRORES CENTINELA (Sentinel Errors)
// Permiten comparación con errors.Is() para manejo específico.
// ============================================================

var (
	// ErrValidacion se produce cuando los datos de entrada no cumplen las reglas de negocio.
	ErrValidacion = errors.New("error de validación")

	// ErrNoEncontrado se produce cuando no se encuentra el recurso solicitado.
	ErrNoEncontrado = errors.New("recurso no encontrado")

	// ErrDuplicado se produce cuando se intenta crear un recurso que ya existe.
	ErrDuplicado = errors.New("recurso duplicado")

	// ErrOperacionInvalida se produce cuando una operación no es permitida en el estado actual.
	ErrOperacionInvalida = errors.New("operación no válida")

	// ErrConexionDB se produce cuando falla la conexión a la base de datos.
	ErrConexionDB = errors.New("error de conexión a base de datos")

	// ErrLimitePrestamos se produce cuando un usuario excede el máximo de préstamos activos.
	ErrLimitePrestamos = errors.New("límite de préstamos alcanzado")

	// ErrLibroNoDisponible se produce cuando se intenta prestar un libro que no está disponible.
	ErrLibroNoDisponible = errors.New("libro no disponible para préstamo")
)

// ============================================================
// TIPO DE ERROR PERSONALIZADO
// Implementa la interfaz error con información adicional de contexto.
// ============================================================

// AppError es un tipo de error personalizado que agrega contexto
// como el código de operación y el recurso afectado.
// Implementa la interfaz error de Go.
type AppError struct {
	Op       string // Operación que falló (ej: "LibroRepository.Crear")
	Recurso  string // Recurso afectado (ej: "libro", "usuario")
	Err      error  // Error original envuelto
}

// Error implementa la interfaz error.
// Retorna un mensaje formateado con la operación y el error original.
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Op, e.Recurso, e.Err)
}

// Unwrap permite a errors.Is() y errors.As() desenvolver el error interno.
// Esto es fundamental para el manejo de errores con cadenas de error en Go.
func (e *AppError) Unwrap() error {
	return e.Err
}

// NuevoAppError crea un error de aplicación con contexto completo.
func NuevoAppError(op, recurso string, err error) *AppError {
	return &AppError{
		Op:      op,
		Recurso: recurso,
		Err:     err,
	}
}
