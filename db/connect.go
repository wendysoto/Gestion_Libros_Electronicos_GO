// Package db gestiona la conexión a la base de datos PostgreSQL.
// Utiliza GORM como ORM y godotenv para cargar variables de entorno.
package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// instancia almacena la conexión singleton a la base de datos.
var instancia *gorm.DB

// LibroDB es el modelo de GORM para la tabla 'libros'.
// Mapea directamente a la tabla de la base de datos PostgreSQL.
type LibroDB struct {
	ID         uint   `gorm:"primaryKey;column:id"`
	Titulo     string `gorm:"column:titulo;not null"`
	Autor      string `gorm:"column:autor;not null"`
	Categoria  string `gorm:"column:categoria;not null"`
	ISBN       string `gorm:"column:isbn;uniqueIndex;not null"`
	Formato    string `gorm:"column:formato;not null;default:'PDF'"`
	Disponible bool   `gorm:"column:disponible;not null;default:true"`
	CreatedAt  string `gorm:"column:created_at"`
	UpdatedAt  string `gorm:"column:updated_at"`
}

// TableName especifica el nombre de la tabla para GORM.
func (LibroDB) TableName() string { return "libros" }

// Connect establece la conexión con la base de datos PostgreSQL.
// Carga las credenciales desde el archivo .env y retorna una instancia de GORM.
// Implementa manejo de errores robusto con mensajes descriptivos.
func Connect() (*gorm.DB, error) {
	// Si ya existe una conexión, reutilizarla (patrón Singleton)
	if instancia != nil {
		return instancia, nil
	}

	// Cargar variables de entorno desde archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar archivo .env, usando variables del sistema")
	}

	// Construir cadena de conexión DSN (Data Source Name) para PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Establecer conexión usando GORM con el driver de PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
	}

	log.Println("✓ Conexión exitosa a la base de datos PostgreSQL")
	instancia = db
	return instancia, nil
}

// GetDB retorna la instancia actual de la base de datos.
// Retorna error si no se ha establecido conexión previamente.
func GetDB() (*gorm.DB, error) {
	if instancia == nil {
		return Connect()
	}
	return instancia, nil
}
