// Sistema de Gestión de Libros Electrónicos
// Proyecto de Programación Orientada a Objetos - Golang
//
// Este archivo es el punto de entrada de la aplicación.
// Configura la conexión a la base de datos, inicializa los servicios
// (inyectando dependencias mediante interfaces) y levanta el servidor HTTP.
//
// Conceptos POO demostrados:
// - Encapsulación: campos privados con getters/setters en los modelos
// - Interfaces: contratos para repositorios y servicios
// - Polimorfismo: múltiples implementaciones de las mismas interfaces
// - Manejo de errores: errores personalizados con wrapping
package main

import (
	"html/template"
	"log"
	"net/http"

	"gestion_libros/db"
	"gestion_libros/handlers"
	"gestion_libros/repositories"
	"gestion_libros/services"

	"github.com/gorilla/mux"
)

func main() {
	// ============================================================
	// PASO 1: Conexión a la Base de Datos PostgreSQL
	// ============================================================
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Error fatal al conectar con la base de datos:", err)
	}
	log.Println("Base de datos conectada exitosamente")

	// ============================================================
	// PASO 2: Inicializar Repositorios (capa de datos)
	// ============================================================
	libroRepo := repositories.NuevoLibroRepository(database)

	// ============================================================
	// PASO 3: Inicializar Servicios (lógica de negocio)
	// ============================================================
	libroService := services.NuevoLibroService(libroRepo)

	// ============================================================
	// PASO 4: Cargar Templates HTML
	// ============================================================
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// ============================================================
	// PASO 5: Inicializar Handlers (controladores HTTP)
	// ============================================================
	libroHandler := handlers.NuevoLibroHandler(libroService, tmpl)

	// ============================================================
	// PASO 6: Configurar Rutas con Gorilla Mux
	// ============================================================
	router := mux.NewRouter()

	// Archivos estáticos (CSS, imágenes)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// Ruta principal
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{"Titulo": "Inicio"}
		tmpl.ExecuteTemplate(w, "index.html", data)
	}).Methods("GET")

	// Rutas de Libros
	router.HandleFunc("/libros", libroHandler.ListarLibros).Methods("GET")
	router.HandleFunc("/libros/nuevo", libroHandler.CrearLibroForm).Methods("GET")
	router.HandleFunc("/libros/nuevo", libroHandler.CrearLibro).Methods("POST")
	router.HandleFunc("/libros/eliminar/{id:[0-9]+}", libroHandler.EliminarLibro).Methods("POST")

	// ============================================================
	// PASO 7: Iniciar Servidor HTTP
	// ============================================================
	puerto := ":8081"
	log.Printf("Servidor iniciado en http://localhost%s", puerto)
	log.Fatal(http.ListenAndServe(puerto, router))
}
