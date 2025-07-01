package repository

import "go-book-clean-architecture-api/internal/domain"

// BookRepository define el contrato para las operaciones de persistencia de libros
// En Clean Architecture, las interfaces están en la capa de dominio/use cases
// y las implementaciones en la capa de infraestructura
// Esto permite la inversión de dependencias: el dominio no depende de la infraestructura
type BookRepository interface {
	// Create almacena un nuevo libro y retorna el libro creado o un error
	Create(book *domain.Book) (*domain.Book, error)

	// GetByID busca un libro por su ID único
	GetByID(id string) (*domain.Book, error)

	// GetAll retorna todos los libros disponibles
	GetAll() ([]*domain.Book, error)

	// Update modifica un libro existente
	Update(book *domain.Book) (*domain.Book, error)

	// Delete elimina un libro por su ID
	Delete(id string) error
}

// UserRepository define el contrato para las operaciones de persistencia de usuarios
// Mantenemos las interfaces simples y enfocadas en operaciones específicas
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
