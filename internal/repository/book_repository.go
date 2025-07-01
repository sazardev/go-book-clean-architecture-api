// Package repository define los CONTRATOS (interfaces) para la persistencia de datos
// Esta es la CAPA DE APLICACIÓN/CASOS DE USO de Clean Architecture
//
// 🎯 CONCEPTO CLAVE: Dependency Inversion Principle (DIP)
// ✅ Las interfaces están aquí (capa de aplicación)
// ✅ Las implementaciones están en infrastructure (capa externa)
// ✅ Esto permite que el dominio NO dependa de detalles técnicos
//
// 🤔 ¿Por qué interfaces?
// - Permiten cambiar implementaciones sin romper código
// - Facilitan el testing (podemos crear mocks fácilmente)
// - Implementan el principio de inversión de dependencias
// - Desacoplan la lógica de negocio de los detalles técnicos
//
// 💡 REGLA DE ORO: "Depend on abstractions, not concretions"
package repository

import "go-book-clean-architecture-api/internal/domain"

// BookRepository define el CONTRATO para las operaciones de persistencia de libros
//
// 📋 ¿Qué es un contrato en programación?
// - Es una interfaz que define QUÉ se puede hacer, pero NO CÓMO hacerlo
// - Cualquier implementación debe cumplir este contrato
// - Permite múltiples implementaciones: memoria, PostgreSQL, MongoDB, etc.
//
// 🔄 Flujo de dependencias:
// UseCase → Repository Interface ← Repository Implementation
//
// 🎯 Beneficios:
// - El caso de uso NO conoce si usamos memoria, PostgreSQL, etc.
// - Podemos cambiar la implementación sin tocar la lógica de negocio
// - Podemos testear fácilmente usando mocks
type BookRepository interface {
	// Create almacena un nuevo libro y retorna el libro creado o un error
	// 📝 Nota: Recibe una entidad completa, no campos separados
	Create(book *domain.Book) (*domain.Book, error)

	// GetByID busca un libro por su ID único
	// 🔍 Retorna error si el libro no existe
	GetByID(id string) (*domain.Book, error)

	// GetAll retorna todos los libros disponibles
	// 📚 En aplicaciones reales, implementarías paginación aquí
	GetAll() ([]*domain.Book, error)

	// Update modifica un libro existente
	// ✏️ Debe verificar que el libro existe antes de actualizar
	Update(book *domain.Book) (*domain.Book, error)

	// Delete elimina un libro por su ID
	// 🗑️ Retorna error si el libro no existe
	Delete(id string) error
}

// UserRepository define el contrato para las operaciones de persistencia de usuarios
//
// 👤 ¿Por qué separamos BookRepository y UserRepository?
// - Principio de Responsabilidad Única (SRP)
// - Cada interfaz tiene una responsabilidad específica
// - Más fácil de mantener y extender
// - Permite implementaciones independientes
//
// 🔧 Nota: En aplicaciones más grandes, podrías tener:
// - BookRepository, UserRepository, OrderRepository, etc.
// - Cada uno enfocado en una entidad específica
type UserRepository interface {
	// Create almacena un nuevo usuario y retorna el usuario creado o un error
	Create(user *domain.User) (*domain.User, error)

	// GetByID busca un usuario por su ID único
	GetByID(id string) (*domain.User, error)

	// GetAll retorna todos los usuarios disponibles
	GetAll() ([]*domain.User, error)

	// Update modifica un usuario existente
	Update(user *domain.User) (*domain.User, error)

	// Delete elimina un usuario por su ID
	Delete(id string) error
}

// 💡 CONSEJOS PARA PRINCIPIANTES:
//
// 1. 🎯 KISS (Keep It Simple, Stupid): Mantén las interfaces simples
// 2. 📝 Usa nombres descriptivos: Create, GetByID, GetAll, Update, Delete
// 3. 🚫 NO incluyas lógica de negocio aquí, solo operaciones de persistencia
// 4. ✅ SÍ piensa en las operaciones que realmente necesitas
//
// 🌟 EJEMPLOS DE MÉTODOS QUE PODRÍAS AGREGAR:
// - GetByAuthor(author string) ([]*domain.Book, error)
// - GetByTitle(title string) (*domain.Book, error)
// - GetByEmailAddress(email string) (*domain.User, error)
// - CountBooks() (int, error)
//
// 🚫 EJEMPLOS DE LO QUE NO DEBES PONER AQUÍ:
// - Validaciones de negocio (van en los casos de uso)
// - Lógica HTTP (va en los handlers)
// - Detalles de base de datos específicos (van en las implementaciones)
//
// 🎓 PREGUNTA PARA REFLEXIONAR:
// ¿Qué pasaría si quisiéramos cambiar de una base de datos en memoria
// a PostgreSQL? ¿Tendríamos que cambiar nuestros casos de uso?
// Respuesta: ¡NO! Solo cambiaríamos la implementación, no el contrato.
