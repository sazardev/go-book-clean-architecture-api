-- Script de inicialización para PostgreSQL
-- Este archivo se ejecuta automáticamente cuando se crea el contenedor de PostgreSQL
-- Útil para cuando quieras implementar un repositorio real con PostgreSQL

-- Crear la base de datos principal (si no existe)
CREATE DATABASE IF NOT EXISTS cleanarch;

-- Conectar a la base de datos
\c cleanarch;

-- Tabla para almacenar libros
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla para almacenar usuarios
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para mejorar performance
CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Función para actualizar automáticamente updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers para actualizar automáticamente updated_at
CREATE TRIGGER update_books_updated_at 
    BEFORE UPDATE ON books 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Datos de ejemplo (opcional)
INSERT INTO books (title, author) VALUES 
    ('Clean Architecture', 'Robert C. Martin'),
    ('The Go Programming Language', 'Alan Donovan'),
    ('Design Patterns', 'Gang of Four')
ON CONFLICT DO NOTHING;

INSERT INTO users (name, email) VALUES 
    ('Juan Pérez', 'juan@example.com'),
    ('María García', 'maria@example.com'),
    ('Carlos López', 'carlos@example.com')
ON CONFLICT (email) DO NOTHING;

-- Mostrar información final
SELECT 'Inicialización de base de datos completada!' as status;
SELECT COUNT(*) as total_books FROM books;
SELECT COUNT(*) as total_users FROM users;
