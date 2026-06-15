// Package models define las estructuras de datos del sistema.
// Implementa encapsulación mediante campos no exportados y métodos getter/setter.
package models

import (
	"fmt"
	"strings"
	"time"
)

// FormatoLibro representa los formatos válidos de un libro electrónico.
type FormatoLibro string

const (
	FormatoPDF  FormatoLibro = "PDF"
	FormatoEPUB FormatoLibro = "EPUB"
	FormatoMOBI FormatoLibro = "MOBI"
)

// Libro representa un libro electrónico en el sistema.
// Los campos están encapsulados (no exportados) para controlar el acceso
// y validar datos mediante métodos setter.
type Libro struct {
	id         uint
	titulo     string
	autor      string
	categoria  string
	isbn       string
	formato    FormatoLibro
	disponible bool
	createdAt  time.Time
	updatedAt  time.Time
}

// --- Métodos Getter (acceso controlado a campos encapsulados) ---

func (l *Libro) GetID() uint            { return l.id }
func (l *Libro) GetTitulo() string       { return l.titulo }
func (l *Libro) GetAutor() string        { return l.autor }
func (l *Libro) GetCategoria() string    { return l.categoria }
func (l *Libro) GetISBN() string         { return l.isbn }
func (l *Libro) GetFormato() FormatoLibro { return l.formato }
func (l *Libro) GetDisponible() bool     { return l.disponible }
func (l *Libro) GetCreatedAt() time.Time { return l.createdAt }
func (l *Libro) GetUpdatedAt() time.Time { return l.updatedAt }

// --- Métodos Setter (validan datos antes de asignarlos) ---

// SetTitulo establece el título del libro.
// Retorna error si el título está vacío o excede 255 caracteres.
func (l *Libro) SetTitulo(titulo string) error {
	titulo = strings.TrimSpace(titulo)
	if titulo == "" {
		return fmt.Errorf("%w: el título no puede estar vacío", ErrValidacion)
	}
	if len(titulo) > 255 {
		return fmt.Errorf("%w: el título no puede exceder 255 caracteres", ErrValidacion)
	}
	l.titulo = titulo
	l.updatedAt = time.Now()
	return nil
}

// SetAutor establece el autor del libro.
// Retorna error si el autor está vacío o excede 255 caracteres.
func (l *Libro) SetAutor(autor string) error {
	autor = strings.TrimSpace(autor)
	if autor == "" {
		return fmt.Errorf("%w: el autor no puede estar vacío", ErrValidacion)
	}
	if len(autor) > 255 {
		return fmt.Errorf("%w: el autor no puede exceder 255 caracteres", ErrValidacion)
	}
	l.autor = autor
	l.updatedAt = time.Now()
	return nil
}

// SetCategoria establece la categoría del libro.
// Retorna error si la categoría está vacía.
func (l *Libro) SetCategoria(categoria string) error {
	categoria = strings.TrimSpace(categoria)
	if categoria == "" {
		return fmt.Errorf("%w: la categoría no puede estar vacía", ErrValidacion)
	}
	l.categoria = categoria
	l.updatedAt = time.Now()
	return nil
}

// SetISBN establece el ISBN del libro.
// Retorna error si el ISBN está vacío o no tiene formato válido.
func (l *Libro) SetISBN(isbn string) error {
	isbn = strings.TrimSpace(isbn)
	if isbn == "" {
		return fmt.Errorf("%w: el ISBN no puede estar vacío", ErrValidacion)
	}
	l.isbn = isbn
	l.updatedAt = time.Now()
	return nil
}

// SetFormato establece el formato del libro electrónico.
// Solo acepta formatos válidos: PDF, EPUB, MOBI.
func (l *Libro) SetFormato(formato FormatoLibro) error {
	switch formato {
	case FormatoPDF, FormatoEPUB, FormatoMOBI:
		l.formato = formato
		l.updatedAt = time.Now()
		return nil
	default:
		return fmt.Errorf("%w: formato '%s' no válido (use PDF, EPUB o MOBI)", ErrValidacion, formato)
	}
}

// SetDisponible cambia el estado de disponibilidad del libro.
func (l *Libro) SetDisponible(disponible bool) {
	l.disponible = disponible
	l.updatedAt = time.Now()
}

// SetID establece el ID del libro (usado internamente por el repositorio).
func (l *Libro) SetID(id uint) {
	l.id = id
}

// SetCreatedAt establece la fecha de creación (uso interno).
func (l *Libro) SetCreatedAt(t time.Time) {
	l.createdAt = t
}

// SetUpdatedAt establece la fecha de actualización (uso interno).
func (l *Libro) SetUpdatedAt(t time.Time) {
	l.updatedAt = t
}

// NuevoLibro es un constructor que crea un libro con validación completa.
// Aplica encapsulación al obligar el uso de setters para la inicialización.
func NuevoLibro(titulo, autor, categoria, isbn string, formato FormatoLibro) (*Libro, error) {
	libro := &Libro{
		disponible: true,
		createdAt:  time.Now(),
		updatedAt:  time.Now(),
	}

	if err := libro.SetTitulo(titulo); err != nil {
		return nil, err
	}
	if err := libro.SetAutor(autor); err != nil {
		return nil, err
	}
	if err := libro.SetCategoria(categoria); err != nil {
		return nil, err
	}
	if err := libro.SetISBN(isbn); err != nil {
		return nil, err
	}
	if err := libro.SetFormato(formato); err != nil {
		return nil, err
	}

	return libro, nil
}

// String implementa la interfaz fmt.Stringer para representación legible.
func (l *Libro) String() string {
	estado := "Disponible"
	if !l.disponible {
		estado = "Prestado"
	}
	return fmt.Sprintf("[%d] %s - %s (%s) [%s]", l.id, l.titulo, l.autor, l.formato, estado)
}
