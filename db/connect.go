// Package db gestiona la conexión a la base de datos PostgreSQL
// y define los modelos de GORM que mapean las tablas del sistema.
package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// instancia almacena la conexión singleton a la base de datos.
var instancia *gorm.DB

// ============================================================
// MODELOS DE BASE DE DATOS (GORM)
// Cada struct mapea a una tabla de PostgreSQL.
// ============================================================

// CategoriaDB mapea la tabla 'categorias'.
type CategoriaDB struct {
	ID          uint   `gorm:"primaryKey;column:id"`
	Nombre      string `gorm:"column:nombre;uniqueIndex;not null"`
	Descripcion string `gorm:"column:descripcion"`
	CreatedAt   string `gorm:"column:created_at"`
	UpdatedAt   string `gorm:"column:updated_at"`
}

func (CategoriaDB) TableName() string { return "categorias" }

// LibroDB mapea la tabla 'libros'.
type LibroDB struct {
	ID          uint   `gorm:"primaryKey;column:id"`
	Titulo      string `gorm:"column:titulo;not null"`
	Autor       string `gorm:"column:autor;not null"`
	CategoriaID uint   `gorm:"column:categoria_id;not null"`
	ISBN        string `gorm:"column:isbn;uniqueIndex;not null"`
	Formato     string `gorm:"column:formato;not null;default:'PDF'"`
	Disponible  bool   `gorm:"column:disponible;not null;default:true"`
	CreatedAt   string `gorm:"column:created_at"`
	UpdatedAt   string `gorm:"column:updated_at"`
}

func (LibroDB) TableName() string { return "libros" }

// UsuarioDB mapea la tabla 'usuarios'.
type UsuarioDB struct {
	ID        uint   `gorm:"primaryKey;column:id"`
	Nombre    string `gorm:"column:nombre;not null"`
	Email     string `gorm:"column:email;uniqueIndex;not null"`
	Tipo      string `gorm:"column:tipo;not null;default:'lector'"`
	Prestados int    `gorm:"column:prestados;not null;default:0"`
	CreatedAt string `gorm:"column:created_at"`
	UpdatedAt string `gorm:"column:updated_at"`
}

func (UsuarioDB) TableName() string { return "usuarios" }

// PrestamoDB mapea la tabla 'prestamos'.
type PrestamoDB struct {
	ID              uint       `gorm:"primaryKey;column:id"`
	LibroID         uint       `gorm:"column:libro_id;not null"`
	UsuarioID       uint       `gorm:"column:usuario_id;not null"`
	FechaPrestamo   time.Time  `gorm:"column:fecha_prestamo;not null"`
	FechaDevolucion *time.Time `gorm:"column:fecha_devolucion"`
	Estado          string     `gorm:"column:estado;not null;default:'activo'"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
}

func (PrestamoDB) TableName() string { return "prestamos" }

// ============================================================
// CONEXION A LA BASE DE DATOS
// ============================================================

// Connect establece la conexión con la base de datos PostgreSQL.
// Carga las credenciales desde el archivo .env y retorna una instancia de GORM.
// Implementa el patrón Singleton para reutilizar la conexión.
func Connect() (*gorm.DB, error) {
	if instancia != nil {
		return instancia, nil
	}

	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar archivo .env, usando variables del sistema")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: false})
	if err != nil {
		return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
	}

	log.Println("Conexion exitosa a la base de datos PostgreSQL")
	instancia = database
	return instancia, nil
}

// GetDB retorna la instancia actual de la base de datos.
func GetDB() (*gorm.DB, error) {
	if instancia == nil {
		return Connect()
	}
	return instancia, nil
}
