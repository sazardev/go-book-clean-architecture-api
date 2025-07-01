package usecase

import (
	"errors"
	"go-book-clean-architecture-api/internal/domain"
	"go-book-clean-architecture-api/internal/repository"
	"strconv"
	"time"
)

// generateID genera un ID único simple usando timestamp
// En un proyecto real, usarías uuid.New().String() o similar
func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// BookUseCase contiene toda la lógica de negocio relacionada con los libros
// Los casos de uso son el corazón de Clean Architecture
// Aquí es donde vive la lógica de negocio específica de la aplicación
type BookUseCase struct {
	bookRepo repository.BookRepository // Dependencia inyectada del repositorio
}

// NewBookUseCase es el constructor que inyecta las dependencias necesarias
// Esta es la implementación del patrón de Inyección de Dependencias
func NewBookUseCase(bookRepo repository.BookRepository) *BookUseCase {
	return &BookUseCase{
		bookRepo: bookRepo,
	}
}

// CreateBook maneja toda la lógica para crear un nuevo libro
// Aquí podríamos agregar validaciones, reglas de negocio, etc.
func (uc *BookUseCase) CreateBook(title, author string) (*domain.Book, error) {
	// Validación de reglas de negocio
	if title == "" {
		return nil, errors.New("el título del libro es obligatorio")
	}
	if author == "" {
		return nil, errors.New("el autor del libro es obligatorio")
	}

	// Crear la entidad del dominio
	book := &domain.Book{
		ID:     generateID(), // Generar ID único
		Title:  title,
		Author: author,
	}

	// Delegar la persistencia al repositorio
	return uc.bookRepo.Create(book)
}

// GetBookByID obtiene un libro por su ID
func (uc *BookUseCase) GetBookByID(id string) (*domain.Book, error) {
	if id == "" {
		return nil, errors.New("ID del libro es obligatorio")
	}
	return uc.bookRepo.GetByID(id)
}

// GetAllBooks obtiene todos los libros disponibles
func (uc *BookUseCase) GetAllBooks() ([]*domain.Book, error) {
	return uc.bookRepo.GetAll()
}

// UpdateBook actualiza un libro existente
func (uc *BookUseCase) UpdateBook(id, title, author string) (*domain.Book, error) {
	if id == "" {
		return nil, errors.New("ID del libro es obligatorio")
	}
	if title == "" {
		return nil, errors.New("el título del libro es obligatorio")
	}
	if author == "" {
		return nil, errors.New("el autor del libro es obligatorio")
	}

	book := &domain.Book{
		ID:     id,
		Title:  title,
		Author: author,
	}

	return uc.bookRepo.Update(book)
}

// DeleteBook elimina un libro por su ID
func (uc *BookUseCase) DeleteBook(id string) error {
	if id == "" {
		return errors.New("ID del libro es obligatorio")
	}
	return uc.bookRepo.Delete(id)
}

// UserUseCase contiene toda la lógica de negocio relacionada con los usuarios
type UserUseCase struct {
	userRepo repository.UserRepository // Dependencia inyectada del repositorio
}

// NewUserUseCase es el constructor que inyecta las dependencias necesarias
func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// CreateUser maneja toda la lógica para crear un nuevo usuario
func (uc *UserUseCase) CreateUser(name, email string) (*domain.User, error) {
	// Validación de reglas de negocio
	if name == "" {
		return nil, errors.New("el nombre del usuario es obligatorio")
	}
	if email == "" {
		return nil, errors.New("el email del usuario es obligatorio")
	}

	// Crear la entidad del dominio
	user := &domain.User{
		ID:    generateID(), // Generar ID único
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
	if id == "" {
		return nil, errors.New("ID del usuario es obligatorio")
	}
	if name == "" {
		return nil, errors.New("el nombre del usuario es obligatorio")
	}
	if email == "" {
		return nil, errors.New("el email del usuario es obligatorio")
	}

	user := &domain.User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	return uc.userRepo.Update(user)
}

// DeleteUser elimina un usuario por su ID
func (uc *UserUseCase) DeleteUser(id string) error {
	if id == "" {
		return errors.New("ID del usuario es obligatorio")
	}
	return uc.userRepo.Delete(id)
}
