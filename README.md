# Biblioteca Digital - Sistema de Gestión de Libros Electrónicos
 
## Descripción
 
Sistema de gestión de libros electrónicos desarrollado en ** Golang** como proyecto de la materia de Programación Orientada a Objetos. Permite administrar un catálogo de libros digitales, gestionar usuarios y controlar el préstamo y devolución de libros mediante una aplicación de consola interactiva.
 
## Objetivo
 
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
biblioteca-digital-go/
├── main.go                  # Punto de entrada y menú principal
├── README.md                # Este archivo
├── .gitignore               # Archivos ignorados por Git
├── libros/
│   └── libros.go           # Lógica del módulo de libros
├── usuarios/
│   └── usuarios.go         # Lógica del módulo de usuarios
├── prestamos/
│   └── prestamos.go        # Lógica del módulo de préstamos
├── utils/
│   └── utils.go            # Funciones auxiliares
└── datos/
    ├── libros.json          # Datos persistidos de libros
    ├── usuarios.json        # Datos persistidos de usuarios
    └── prestamos.json       # Datos persistidos de préstamos
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
 
## Conceptos de Go Aplicados (Unidad 1)
 
| Tema | Aplicación |
|------|-----------|
| Sintaxis y tipos de datos | Structs para modelar Libro, Usuario y Prestamo |
| Condicionales (if/else, switch) | Validaciones de negocio y menú de opciones |
| Bucles (for, range) | Recorrido de catálogos, búsquedas y listados |
| Funciones | Operaciones CRUD con retorno múltiple |
| Paquetes | Organización modular del código por dominio |

## Autores
Estudiantes de la carrera de Ingenieria en Sistemas de Información 

**- Wendy Soto**
**- Jhoel Amagua**
