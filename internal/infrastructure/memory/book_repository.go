package memory

import (
	"errors"
	"go-book-clean-architecture-api/internal/domain"
	"go-book-clean-architecture-api/internal/repository"
	"sync"
)

// InMemoryBookRepository es una implementación en memoria del BookRepository
// Esta implementación está en la capa de infraestructura
// En un caso real, aquí tendríamos implementaciones para PostgreSQL, MongoDB, etc.
type InMemoryBookRepository struct {
	books map[string]*domain.Book // Almacenamiento en memoria usando un map
	mutex sync.RWMutex            // Para manejar concurrencia de manera segura
}

// NewInMemoryBookRepository crea una nueva instancia del repositorio en memoria
func NewInMemoryBookRepository() repository.BookRepository {
	return &InMemoryBookRepository{
		books: make(map[string]*domain.Book),
		mutex: sync.RWMutex{},
	}
}

// Create almacena un nuevo libro en memoria
func (r *InMemoryBookRepository) Create(book *domain.Book) (*domain.Book, error) {
	r.mutex.Lock()         // Bloquear para escritura
	defer r.mutex.Unlock() // Asegurar que se desbloquee al final

	// Verificar si el libro ya existe
	if _, exists := r.books[book.ID]; exists {
		return nil, errors.New("el libro con este ID ya existe")
	}

	// Almacenar el libro
	r.books[book.ID] = book
	return book, nil
}

// GetByID busca un libro por su ID
func (r *InMemoryBookRepository) GetByID(id string) (*domain.Book, error) {
	r.mutex.RLock()         // Bloquear solo para lectura
	defer r.mutex.RUnlock() // Asegurar que se desbloquee al final

	book, exists := r.books[id]
	if !exists {
		return nil, errors.New("libro no encontrado")
	}

	return book, nil
}

// GetAll retorna todos los libros almacenados
func (r *InMemoryBookRepository) GetAll() ([]*domain.Book, error) {
	r.mutex.RLock()         // Bloquear solo para lectura
	defer r.mutex.RUnlock() // Asegurar que se desbloquee al final

	books := make([]*domain.Book, 0, len(r.books))
	for _, book := range r.books {
		books = append(books, book)
	}

	return books, nil
}

// Update modifica un libro existente
func (r *InMemoryBookRepository) Update(book *domain.Book) (*domain.Book, error) {
	r.mutex.Lock()         // Bloquear para escritura
	defer r.mutex.Unlock() // Asegurar que se desbloquee al final

	// Verificar si el libro existe
	if _, exists := r.books[book.ID]; !exists {
		return nil, errors.New("libro no encontrado")
	}

	// Actualizar el libro
	r.books[book.ID] = book
	return book, nil
}

// Delete elimina un libro por su ID
func (r *InMemoryBookRepository) Delete(id string) error {
	r.mutex.Lock()         // Bloquear para escritura
	defer r.mutex.Unlock() // Asegurar que se desbloquee al final

	// Verificar si el libro existe
	if _, exists := r.books[id]; !exists {
		return errors.New("libro no encontrado")
	}

	// Eliminar el libro
	delete(r.books, id)
	return nil
}

// InMemoryUserRepository es una implementación en memoria del UserRepository
type InMemoryUserRepository struct {
	users map[string]*domain.User // Almacenamiento en memoria usando un map
	mutex sync.RWMutex            // Para manejar concurrencia de manera segura
}

// NewInMemoryUserRepository crea una nueva instancia del repositorio en memoria
func NewInMemoryUserRepository() repository.UserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
		mutex: sync.RWMutex{},
	}
}

// Create almacena un nuevo usuario en memoria
func (r *InMemoryUserRepository) Create(user *domain.User) (*domain.User, error) {
	r.mutex.Lock()         // Bloquear para escritura
	defer r.mutex.Unlock() // Asegurar que se desbloquee al final

	// Verificar si el usuario ya existe
	if _, exists := r.users[user.ID]; exists {
		return nil, errors.New("el usuario con este ID ya existe")
	}

	// Almacenar el usuario
	r.users[user.ID] = user
	return user, nil
}

// GetByID busca un usuario por su ID
func (r *InMemoryUserRepository) GetByID(id string) (*domain.User, error) {
	r.mutex.RLock()         // Bloquear solo para lectura
	defer r.mutex.RUnlock() // Asegurar que se desbloquee al final

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("usuario no encontrado")
	}

	return user, nil
}

// GetAll retorna todos los usuarios almacenados
func (r *InMemoryUserRepository) GetAll() ([]*domain.User, error) {
	r.mutex.RLock()         // Bloquear solo para lectura
	defer r.mutex.RUnlock() // Asegurar que se desbloquee al final

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

// Update modifica un usuario existente
func (r *InMemoryUserRepository) Update(user *domain.User) (*domain.User, error) {
	r.mutex.Lock()         // Bloquear para escritura
	defer r.mutex.Unlock() // Asegurar que se desbloquee al final

	// Verificar si el usuario existe
	if _, exists := r.users[user.ID]; !exists {
		return nil, errors.New("usuario no encontrado")
	}

	// Actualizar el usuario
	r.users[user.ID] = user
	return user, nil
}

// Delete elimina un usuario por su ID
func (r *InMemoryUserRepository) Delete(id string) error {
	r.mutex.Lock()         // Bloquear para escritura
	defer r.mutex.Unlock() // Asegurar que se desbloquee al final

	// Verificar si el usuario existe
	if _, exists := r.users[id]; !exists {
		return errors.New("usuario no encontrado")
	}

	// Eliminar el usuario
	delete(r.users, id)
	return nil
}
