-- ============================================================
-- Se crea un primer esquema borrador de base de datos para un sistema de gestión de libros electrónicos en una biblioteca digital.
-- Base de datos compatible: PostgreSQL 14+

-- ============================================================
-- TIPOS ENUMERADOS
-- ============================================================
CREATE TYPE formato_libro AS ENUM ('PDF', 'EPUB', 'MOBI');
CREATE TYPE tipo_usuario AS ENUM ('admin', 'lector');
CREATE TYPE estado_prestamo AS ENUM ('activo', 'devuelto');
 
-- ============================================================
-- TABLA: LIBROS
-- Almacena el catálogo de libros electrónicos
-- ============================================================
CREATE TABLE libros (
    id          SERIAL PRIMARY KEY,
    titulo      VARCHAR(255) NOT NULL,
    autor       VARCHAR(255) NOT NULL,
    categoria   VARCHAR(100) NOT NULL,
    isbn        VARCHAR(20)  UNIQUE NOT NULL,
    formato     formato_libro NOT NULL DEFAULT 'PDF',
    disponible  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
 
-- ============================================================
-- TABLA: USUARIOS
-- Almacena los usuarios del sistema
-- ============================================================
CREATE TABLE usuarios (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(255) NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
    tipo        tipo_usuario NOT NULL DEFAULT 'lector',
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
 
-- ============================================================
-- TABLA: PRESTAMOS
-- Registra los préstamos y devoluciones de libros
-- ============================================================
CREATE TABLE prestamos (
    id                SERIAL PRIMARY KEY,
    libro_id          INT NOT NULL,
    usuario_id        INT NOT NULL,
    fecha_prestamo    DATE NOT NULL,
    fecha_devolucion  DATE NULL,
    estado            estado_prestamo NOT NULL DEFAULT 'activo',
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_libro
        FOREIGN KEY (libro_id) REFERENCES libros(id) ON DELETE RESTRICT,
    CONSTRAINT fk_usuario
        FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE RESTRICT
);
 
-- ============================================================
-- ÍNDICES para mejorar rendimiento de consultas
-- ============================================================
CREATE INDEX idx_libros_categoria ON libros(categoria);
CREATE INDEX idx_libros_autor ON libros(autor);
CREATE INDEX idx_libros_disponible ON libros(disponible);
CREATE INDEX idx_usuarios_tipo ON usuarios(tipo);
CREATE INDEX idx_prestamos_estado ON prestamos(estado);
CREATE INDEX idx_prestamos_usuario ON prestamos(usuario_id, estado);
CREATE INDEX idx_prestamos_libro ON prestamos(libro_id, estado);
 

-- ============================================================
-- DATOS DE PRUEBA
-- ============================================================
 
-- Insertar usuarios
INSERT INTO usuarios (nombre, email, tipo) VALUES
('Administrador', 'admin@biblioteca.com', 'admin'),
('María García', 'maria.garcia@email.com', 'lector'),
('Carlos López', 'carlos.lopez@email.com', 'lector'),
('Ana Martínez', 'ana.martinez@email.com', 'lector');
 
-- Insertar libros
INSERT INTO libros (titulo, autor, categoria, isbn, formato, disponible) VALUES
('Clean Code', 'Robert C. Martin', 'Programación', '978-0132350884', 'PDF', TRUE),
('El Principito', 'Antoine de Saint-Exupéry', 'Literatura', '978-0156012195', 'EPUB', TRUE),
('Introducción a los Algoritmos', 'Thomas H. Cormen', 'Programación', '978-0262033848', 'PDF', FALSE),
('Cien Años de Soledad', 'Gabriel García Márquez', 'Literatura', '978-0307474728', 'MOBI', TRUE),
('The Go Programming Language', 'Alan Donovan', 'Programación', '978-0134190440', 'PDF', TRUE),
('Don Quijote de la Mancha', 'Miguel de Cervantes', 'Literatura', '978-8420412146', 'EPUB', TRUE),
('Design Patterns', 'Erich Gamma', 'Programación', '978-0201633610', 'PDF', TRUE),
('Harry Potter y la Piedra Filosofal', 'J.K. Rowling', 'Ficción', '978-8498386318', 'EPUB', FALSE);
 
-- Insertar préstamos
INSERT INTO prestamos (libro_id, usuario_id, fecha_prestamo, fecha_devolucion, estado) VALUES
(3, 2, '2026-05-10', NULL, 'activo'),
(8, 2, '2026-05-15', NULL, 'activo'),
(1, 3, '2026-05-01', '2026-05-12', 'devuelto'),
(5, 4, '2026-05-20', NULL, 'activo');
 
