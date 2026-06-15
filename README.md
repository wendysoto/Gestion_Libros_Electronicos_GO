## Sistema de Gestión de Libros Electrónicos :closed_book:
 
### Descripción
 
Sistema de gestión de libros electrónicos desarrollado en ** Golang** como proyecto de la materia de Programación Orientada a Objetos. Permite administrar un catálogo de libros digitales, gestionar usuarios y controlar el préstamo y devolución de libros mediante una aplicación de consola interactiva.
 
### Objetivo
 
Desarrollar un sistema funcional que permita la gestión  de una biblioteca digital, aplicando los conceptos fundamentales de programación en Go: sintaxis básica, condicionales, estructuras de control iterativo, funciones y manejo de paquetes.
 
## Módulos
 
### Módulo 1: Gestión de Libros (`libros/`)
 
Administra el catálogo de libros electrónicos disponibles en la biblioteca.
 
**Funcionalidades:**
- Registrar nuevo libro 
- Listar todos los libros del catálogo
- Buscar libros por título, autor o categoría
- Actualizar información de un libro existente
- Eliminar un libro del catálogo
- Consultar disponibilidad
 
### Módulo 2: Gestión de Usuarios (`usuarios/`)
 
Administra los usuarios registrados en el sistema.
 
**Funcionalidades:**
- Registrar nuevo usuario (nombre, email, tipo)
- Listar todos los usuarios
- Buscar usuario por ID o email
- Actualizar datos de un usuario
- Eliminar un usuario del sistema
 
### Módulo 3: Gestión de Préstamos (`prestamos/`)
 
Controla el flujo de préstamos y devoluciones de libros.
 
**Funcionalidades:**
- Administrar préstamos
- Crear un préstamo (asignar libro a usuario)
- Validar disponibilidad del libro antes del préstamo
- Controlar límite de préstamos por usuario (máximo 3 activos)
- Registrar devolución de un libro
- Listar préstamos activos
- Consultar historial de préstamos por usuario
 
### Módulo 4: Validaciones (`Utils/`)
 
Funciones auxiliares compartidas entre los demás módulos.
 
**Funcionalidades:**
- Validación de datos de entrada
- Lectura y escritura de archivos JSON (persistencia)
- Generación de IDs únicos
- Formateo de fechas
 
## Estructura del Proyecto
 
```
Gestion_Libros_Electronicos_GO/
├── README.md                # Este archivo
├── .gitignore               # Archivos ignorados por Git
├── datos/
│   └── base_datos.sql           # Lógica del módulo de libros
├── db/
│   └── connect.go         # Conexión con la base de datos en Postgresql
├── handlers/
│   └── libro_handler.go        # Peticiones http y renderiza vistas en html - encapsulación y manejo de errores de libros
├── interfaces/
│   └── repository.go            # Conjunto de métodos
└── models/
    ├── errors.go          # Errores personalizados
    ├── libro.go           # Estructura - Getter y Setter de la lógica de libros
└── prestamos/
    ├── prestamos.go          # En proceso de desarrollo
└── repositories/
    ├── libro_repository.go          # Lógica de repositorios de libros
└── services/
    ├── libro_service.go          # intermediarios entre los handlers http y los repositorios (base de datos), mantiene la logica del negocio
└── static/
    ├── style.css          # estilos básicos (colores, magenes, etc.)
└── templates/
    ├── base.html          # Encabezado - Menú - Pie de página
    ├── index.html         # Encabezado - Página de inicio - Pie de página
    ├── libro_form.html    # Encabezado - Formulario para ingreso de nuevo libro - Pie de página
    ├── libros.html        # Encabezado - Listado de libros desde base de datos - Pie de página
```
 
## Tecnologías Utilizadas
 
- **Lenguaje:** Go 
- **Persistencia:** Archivos JSON
- **Interfaz:**  
 
## Paquetes Utilizados
 
### Biblioteca Estándar de Go
 
| Paquete | Uso |
|---------|-----|
| `fmt` | Entrada y salida por consola |
| `os` | Manejo de archivos para persistencia |
| `encoding/json` | Serialización y deserialización de datos |
| `strings` | Operaciones de búsqueda y comparación de texto |
| `strconv` | Conversión entre tipos de datos |
| `time` | Manejo de fechas para préstamos |
| `errors` | Manejo de errores personalizados |
| `bufio` | Lectura de entrada del usuario |
 
### Paquetes de Terceros (Opcionales)
 
| Paquete | Uso | Instalación |
|---------|-----|-------------|
| `github.com/fatih/color` | Colorear salida en consola | `go get github.com/fatih/color` |
| `github.com/olekukonez/tablewriter` | Tablas formateadas en consola | `go get github.com/olekukonez/tablewriter` |
 
## Cómo Ejecutar
 
```
# Ejecutar el programa
go run main.go
```
 

 ----------------------------------

 Semana 2


 Uso de Bosstrap a travez del CDN
 https://getbootstrap.com/
 

## Autores
Estudiantes de la carrera de Ingenieria en Sistemas de Información 

**- Wendy Soto**
**- Jhoel Amagua**
