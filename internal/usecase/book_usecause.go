// Package usecase contiene la LÓGICA DE NEGOCIO de la aplicación
// Esta es la CAPA DE APLICACIÓN de Clean Architecture
//
// 🧠 ¿Qué son los casos de uso?
// - Son las reglas de negocio específicas de la aplicación
// - Orquestan las operaciones entre entidades
// - Definen QUÉ puede hacer el usuario en la aplicación
// - Son independientes de frameworks, bases de datos, UI, etc.
//
// 🎯 REGLAS DE ORO de los Casos de Uso:
// ✅ Contienen la lógica de negocio específica de la aplicación
// ✅ Validan datos antes de procesarlos
// ✅ Coordinan operaciones entre diferentes entidades
// ✅ Usan las interfaces de repositorio (no implementaciones concretas)
// ✅ Retornan errores de negocio significativos
//
// 🚫 Los casos de uso NO deben:
// - Conocer detalles de HTTP (request/response)
// - Conocer detalles de base de datos (SQL, tablas)
// - Depender de frameworks externos
// - Manejar presentación de datos
package usecase

import (
	"errors"
	"go-book-clean-architecture-api/internal/domain"
	"go-book-clean-architecture-api/internal/repository"

	"github.com/google/uuid"
)

// BookUseCase contiene toda la lógica de negocio relacionada con los libros
//
// 📚 ¿Por qué separamos la lógica en casos de uso?
// - Centraliza las reglas de negocio en un lugar
// - Facilita el testing (sin HTTP, sin DB)
// - Permite reutilizar la lógica desde diferentes interfaces (HTTP, CLI, gRPC)
// - Hace el código más mantenible y escalable
//
// 🔧 Patrón de Inyección de Dependencias:
// - El caso de uso RECIBE las dependencias que necesita
// - NO crea las dependencias internamente
// - Esto facilita el testing y la flexibilidad
type BookUseCase struct {
	bookRepo repository.BookRepository // Dependencia inyectada del repositorio
}

// NewBookUseCase es el CONSTRUCTOR que implementa Dependency Injection
//
// 🎯 ¿Por qué usar constructores?
// - Garantizan que el objeto se crea correctamente
// - Hacen explícitas las dependencias necesarias
// - Facilitan el testing (podemos inyectar mocks)
// - Siguen el principio de inversión de dependencias
//
// 💡 Nota: En Go, los constructores son por convención funciones New*
func NewBookUseCase(bookRepo repository.BookRepository) *BookUseCase {
	return &BookUseCase{
		bookRepo: bookRepo,
	}
}

// CreateBook maneja toda la lógica para crear un nuevo libro
//
// 🔄 Flujo típico de un caso de uso:
// 1. Validar entrada (reglas de negocio)
// 2. Crear/modificar entidades del dominio
// 3. Llamar al repositorio para persistir
// 4. Retornar resultado o error
//
// 🎯 Responsabilidades:
// ✅ Validar que el título no esté vacío (regla de negocio)
// ✅ Validar que el autor no esté vacío (regla de negocio)
// ✅ Generar ID único para el libro
// ✅ Crear la entidad Book
// ✅ Delegar la persistencia al repositorio
func (uc *BookUseCase) CreateBook(title, author string) (*domain.Book, error) {
	// PASO 1: Validaciones de reglas de negocio
	// Estas son reglas específicas de nuestro dominio
	if title == "" {
		return nil, errors.New("el título del libro es obligatorio")
	}
	if author == "" {
		return nil, errors.New("el autor del libro es obligatorio")
	}

	// PASO 2: Crear la entidad del dominio
	book := &domain.Book{
		ID:     uuid.New().String(), // Generar ID único
		Title:  title,
		Author: author,
	}

	// PASO 3: Delegar la persistencia al repositorio
	// El caso de uso NO sabe si esto se guarda en memoria, PostgreSQL, etc.
	return uc.bookRepo.Create(book)
}

// GetBookByID obtiene un libro por su ID
//
// 🔍 Caso de uso simple: validar entrada y delegar al repositorio
// Podríamos agregar lógica adicional como logging, métricas, cache, etc.
func (uc *BookUseCase) GetBookByID(id string) (*domain.Book, error) {
	// Validación de entrada
	if id == "" {
		return nil, errors.New("ID del libro es obligatorio")
	}

	// Delegar al repositorio
	return uc.bookRepo.GetByID(id)
}

// GetAllBooks obtiene todos los libros disponibles
//
// 📚 En aplicaciones reales, aquí podrías implementar:
// - Paginación: GetBooks(page, limit int)
// - Filtros: GetBooksByAuthor(author string)
// - Ordenamiento: GetBooksSortedByTitle()
// - Cache: verificar cache antes de llamar al repositorio
func (uc *BookUseCase) GetAllBooks() ([]*domain.Book, error) {
	return uc.bookRepo.GetAll()
}

// UpdateBook actualiza un libro existente
//
// 🔄 Lógica de actualización:
// 1. Validar todos los campos (título, autor, ID)
// 2. Crear entidad con los nuevos datos
// 3. Delegar la actualización al repositorio
//
// 💡 Nota: El repositorio se encarga de verificar si el libro existe
func (uc *BookUseCase) UpdateBook(id, title, author string) (*domain.Book, error) {
	// Validaciones de negocio
	if id == "" {
		return nil, errors.New("ID del libro es obligatorio")
	}
	if title == "" {
		return nil, errors.New("el título del libro es obligatorio")
	}
	if author == "" {
		return nil, errors.New("el autor del libro es obligatorio")
	}

	// Crear entidad con los datos actualizados
	book := &domain.Book{
		ID:     id,
		Title:  title,
		Author: author,
	}

	// Delegar la actualización al repositorio
	return uc.bookRepo.Update(book)
}

// DeleteBook elimina un libro por su ID
//
// 🗑️ Operación de eliminación simple
// En aplicaciones reales, podrías implementar:
// - Soft delete (marcar como eliminado, no borrar físicamente)
// - Verificaciones adicionales (¿el libro está prestado?)
// - Logging de auditoría
func (uc *BookUseCase) DeleteBook(id string) error {
	// Validación de entrada
	if id == "" {
		return errors.New("ID del libro es obligatorio")
	}

	// Delegar la eliminación al repositorio
	return uc.bookRepo.Delete(id)
}

// UserUseCase contiene toda la lógica de negocio relacionada con los usuarios
//
// 👤 Misma estructura que BookUseCase, pero para usuarios
// Esto demuestra el patrón consistente en Clean Architecture
type UserUseCase struct {
	userRepo repository.UserRepository // Dependencia inyectada del repositorio
}

// NewUserUseCase constructor para UserUseCase
func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// CreateUser maneja toda la lógica para crear un nuevo usuario
//
// 👤 Lógica específica para usuarios:
// - Validar que el nombre no esté vacío
// - Validar que el email no esté vacío
// - En aplicaciones reales: validar formato de email, unicidad, etc.
func (uc *UserUseCase) CreateUser(name, email string) (*domain.User, error) {
	// Validaciones de reglas de negocio
	if name == "" {
		return nil, errors.New("el nombre del usuario es obligatorio")
	}
	if email == "" {
		return nil, errors.New("el email del usuario es obligatorio")
	}

	// TODO: En aplicaciones reales, aquí validarías:
	// - Formato de email válido
	// - Email único en el sistema
	// - Longitud mínima del nombre
	// - Caracteres permitidos, etc.

	// Crear la entidad del dominio
	user := &domain.User{
		ID:    uuid.New().String(), // Generar ID único
		Name:  name,
		Email: email,
	}

	// Delegar la persistencia al repositorio
	return uc.userRepo.Create(user)
}

// GetUserByID obtiene un usuario por su ID
func (uc *UserUseCase) GetUserByID(id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("ID del usuario es obligatorio")
	}
	return uc.userRepo.GetByID(id)
}

// GetAllUsers obtiene todos los usuarios disponibles
func (uc *UserUseCase) GetAllUsers() ([]*domain.User, error) {
	return uc.userRepo.GetAll()
}

// UpdateUser actualiza un usuario existente
func (uc *UserUseCase) UpdateUser(id, name, email string) (*domain.User, error) {
	// Validaciones de negocio
	if id == "" {
		return nil, errors.New("ID del usuario es obligatorio")
	}
	if name == "" {
		return nil, errors.New("el nombre del usuario es obligatorio")
	}
	if email == "" {
		return nil, errors.New("el email del usuario es obligatorio")
	}

	// Crear entidad con los datos actualizados
	user := &domain.User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	// Delegar la actualización al repositorio
	return uc.userRepo.Update(user)
}

// DeleteUser elimina un usuario por su ID
func (uc *UserUseCase) DeleteUser(id string) error {
	if id == "" {
		return errors.New("ID del usuario es obligatorio")
	}
	return uc.userRepo.Delete(id)
}

// 💡 CONSEJOS PARA PRINCIPIANTES:
//
// 1. 🎯 Un caso de uso = Una operación específica del negocio
// 2. 📝 Siempre validar entradas (nunca confíes en datos externos)
// 3. 🧠 Piensa en las reglas de negocio, no en detalles técnicos
// 4. 🔧 Usa inyección de dependencias para mantener flexibilidad
// 5. ✅ Retorna errores descriptivos que aporten valor al usuario
//
// 🌟 EJEMPLOS DE CASOS DE USO ADICIONALES QUE PODRÍAS AGREGAR:
// - SearchBooksByAuthor(author string) ([]*domain.Book, error)
// - GetBookStatistics() (*domain.BookStats, error)
// - LendBookToUser(bookID, userID string) error
// - GetUserBorrowedBooks(userID string) ([]*domain.Book, error)
//
// 🚫 EJEMPLOS DE LO QUE NO DEBES PONER AQUÍ:
// - Detalles de HTTP (parsing JSON, status codes)
// - Detalles de base de datos (SQL queries, conexiones)
// - Lógica de presentación (formateo de fechas, UI)
// - Dependencias de frameworks específicos
//
// 🎓 PREGUNTA PARA REFLEXIONAR:
// Si quisiéramos exponer esta misma lógica a través de una CLI
// en lugar de una API HTTP, ¿tendríamos que cambiar estos casos de uso?
// Respuesta: ¡NO! Solo crearíamos nuevos handlers para CLI.
