// Sistema de Gestión de Libros Electrónicos
// Proyecto de Programación Orientada a Objetos - Golang
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
	// Conexión a la Base de Datos PostgreSQL
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Error fatal al conectar con la base de datos:", err)
	}
	log.Println("Base de datos conectada exitosamente")

	// Inicializar Repositorios
	categoriaRepo := repositories.NuevoCategoriaRepository(database)
	libroRepo := repositories.NuevoLibroRepository(database)
	usuarioRepo := repositories.NuevoUsuarioRepository(database)
	prestamoRepo := repositories.NuevoPrestamoRepository(database)

	// Inicializar Servicios
	categoriaService := services.NuevoCategoriaService(categoriaRepo)
	libroService := services.NuevoLibroService(libroRepo)
	usuarioService := services.NuevoUsuarioService(usuarioRepo, prestamoRepo)
	prestamoService := services.NuevoPrestamoService(prestamoRepo, libroRepo, usuarioRepo)

	// Cargar Templates HTML
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// Inicializar Handlers
	libroHandler := handlers.NuevoLibroHandler(libroService, categoriaService, tmpl)
	usuarioHandler := handlers.NuevoUsuarioHandler(usuarioService, tmpl)
	prestamoHandler := handlers.NuevoPrestamoHandler(prestamoService, libroService, usuarioService, categoriaService, tmpl)
	apiHandler := handlers.NuevoAPIHandler(libroService, categoriaService)

	// Configurar Rutas
	router := mux.NewRouter()

	// Archivos estáticos
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
	router.HandleFunc("/libros/buscar", libroHandler.BuscarLibrosPorCategoria).Methods("GET")
	router.HandleFunc("/libros/nuevo", libroHandler.CrearLibroForm).Methods("GET")
	router.HandleFunc("/libros/nuevo", libroHandler.CrearLibro).Methods("POST")
	router.HandleFunc("/libros/editar/{id:[0-9]+}", libroHandler.EditarLibroForm).Methods("GET")
	router.HandleFunc("/libros/editar/{id:[0-9]+}", libroHandler.EditarLibro).Methods("POST")
	router.HandleFunc("/libros/eliminar/{id:[0-9]+}", libroHandler.EliminarLibro).Methods("POST")

	// Rutas de Usuarios
	router.HandleFunc("/usuarios", usuarioHandler.ListarUsuarios).Methods("GET")
	router.HandleFunc("/usuarios/nuevo", usuarioHandler.CrearUsuarioForm).Methods("GET")
	router.HandleFunc("/usuarios/nuevo", usuarioHandler.CrearUsuario).Methods("POST")
	router.HandleFunc("/usuarios/editar/{id:[0-9]+}", usuarioHandler.EditarUsuarioForm).Methods("GET")
	router.HandleFunc("/usuarios/editar/{id:[0-9]+}", usuarioHandler.EditarUsuario).Methods("POST")
	router.HandleFunc("/usuarios/eliminar/{id:[0-9]+}", usuarioHandler.EliminarUsuario).Methods("POST")

	// Rutas de Préstamos
	router.HandleFunc("/prestamos", prestamoHandler.ListarPrestamos).Methods("GET")
	router.HandleFunc("/prestamos/nuevo", prestamoHandler.CrearPrestamoForm).Methods("GET")
	router.HandleFunc("/prestamos/nuevo", prestamoHandler.CrearPrestamo).Methods("POST")
	router.HandleFunc("/prestamos/devolver/{id:[0-9]+}", prestamoHandler.DevolverPrestamo).Methods("POST")
	router.HandleFunc("/prestamos/usuario/{id:[0-9]+}", prestamoHandler.HistorialPorUsuario).Methods("GET")

	// API REST (JSON)
	router.HandleFunc("/api/libros", apiHandler.ListarLibrosPorCategoria).Methods("GET")

	// Iniciar Servidor
	puerto := ":8081"
	log.Printf("Servidor iniciado en http://localhost%s", puerto)
	log.Fatal(http.ListenAndServe(puerto, router))
}
