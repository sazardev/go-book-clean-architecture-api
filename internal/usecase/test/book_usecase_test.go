// Package test contiene ejemplos de testing para casos de uso
// Esto demuestra cómo testear la lógica de negocio sin dependencias externas
//
// 🧪 ¿Por qué es importante testear los casos de uso?
// - Son el corazón de la lógica de negocio
// - Se pueden testear sin HTTP, base de datos, o frameworks
// - Los tests son rápidos y confiables
// - Nos dan confianza al hacer cambios
//
// 🎯 BENEFICIOS de testear en esta capa:
// ✅ Tests rápidos (sin I/O)
// ✅ Tests independientes (sin estado compartido)
// ✅ Tests enfocados en lógica de negocio
// ✅ Fácil crear diferentes escenarios
package test

import (
	"errors"
	"go-book-clean-architecture-api/internal/domain"
	"go-book-clean-architecture-api/internal/usecase"
	"testing"
)

// MockBookRepository es un mock simple del BookRepository para testing
//
// 🎭 ¿Qué es un mock?
// - Una implementación falsa que controla qué devuelve
// - Permite testear casos específicos (éxito, error, etc.)
// - No tiene efectos secundarios (no guarda en DB real)
// - Es predecible y controlable
type MockBookRepository struct {
	books       map[string]*domain.Book
	shouldError bool
}

// NewMockBookRepository crea una nueva instancia del mock
func NewMockBookRepository() *MockBookRepository {
	return &MockBookRepository{
		books:       make(map[string]*domain.Book),
		shouldError: false,
	}
}

// SetShouldError configura si el mock debe retornar errores
func (m *MockBookRepository) SetShouldError(shouldError bool) {
	m.shouldError = shouldError
}

// Implementación de la interfaz BookRepository

func (m *MockBookRepository) Create(book *domain.Book) (*domain.Book, error) {
	if m.shouldError {
		return nil, errors.New("error simulado del repositorio")
	}
	m.books[book.ID] = book
	return book, nil
}

func (m *MockBookRepository) GetByID(id string) (*domain.Book, error) {
	if m.shouldError {
		return nil, errors.New("error simulado del repositorio")
	}
	book, exists := m.books[id]
	if !exists {
		return nil, errors.New("libro no encontrado")
	}
	return book, nil
}

func (m *MockBookRepository) GetAll() ([]*domain.Book, error) {
	if m.shouldError {
		return nil, errors.New("error simulado del repositorio")
	}
	books := make([]*domain.Book, 0, len(m.books))
	for _, book := range m.books {
		books = append(books, book)
	}
	return books, nil
}

func (m *MockBookRepository) Update(book *domain.Book) (*domain.Book, error) {
	if m.shouldError {
		return nil, errors.New("error simulado del repositorio")
	}
	if _, exists := m.books[book.ID]; !exists {
		return nil, errors.New("libro no encontrado")
	}
	m.books[book.ID] = book
	return book, nil
}

func (m *MockBookRepository) Delete(id string) error {
	if m.shouldError {
		return errors.New("error simulado del repositorio")
	}
	if _, exists := m.books[id]; !exists {
		return errors.New("libro no encontrado")
	}
	delete(m.books, id)
	return nil
}

// TestCreateBook_Success prueba el caso exitoso de crear un libro
//
// 🧪 Patrón AAA (Arrange-Act-Assert):
// - Arrange: Preparar el entorno de testing
// - Act: Ejecutar la acción que queremos testear
// - Assert: Verificar que el resultado es el esperado
func TestCreateBook_Success(t *testing.T) {
	// Arrange: Preparar el entorno
	mockRepo := NewMockBookRepository()
	bookUseCase := usecase.NewBookUseCase(mockRepo)

	// Act: Ejecutar la acción
	book, err := bookUseCase.CreateBook("Clean Architecture", "Robert C. Martin")

	// Assert: Verificar resultados
	if err != nil {
		t.Errorf("Se esperaba que no hubiera error, pero se obtuvo: %v", err)
	}
	if book == nil {
		t.Fatal("Se esperaba un libro, pero se obtuvo nil")
	}
	if book.Title != "Clean Architecture" {
		t.Errorf("Se esperaba título 'Clean Architecture', pero se obtuvo: %s", book.Title)
	}
	if book.Author != "Robert C. Martin" {
		t.Errorf("Se esperaba autor 'Robert C. Martin', pero se obtuvo: %s", book.Author)
	}
	if book.ID == "" {
		t.Error("Se esperaba que el libro tuviera un ID generado")
	}
}

// TestCreateBook_EmptyTitle prueba el error cuando el título está vacío
func TestCreateBook_EmptyTitle(t *testing.T) {
	// Arrange
	mockRepo := NewMockBookRepository()
	bookUseCase := usecase.NewBookUseCase(mockRepo)

	// Act
	book, err := bookUseCase.CreateBook("", "Algún autor")

	// Assert
	if err == nil {
		t.Error("Se esperaba un error, pero no se obtuvo ninguno")
	}
	if book != nil {
		t.Error("Se esperaba nil, pero se obtuvo un libro")
	}
	expectedError := "el título del libro es obligatorio"
	if err.Error() != expectedError {
		t.Errorf("Se esperaba error '%s', pero se obtuvo: %s", expectedError, err.Error())
	}
}

// TestCreateBook_EmptyAuthor prueba el error cuando el autor está vacío
func TestCreateBook_EmptyAuthor(t *testing.T) {
	// Arrange
	mockRepo := NewMockBookRepository()
	bookUseCase := usecase.NewBookUseCase(mockRepo)

	// Act
	book, err := bookUseCase.CreateBook("Algún título", "")

	// Assert
	if err == nil {
		t.Error("Se esperaba un error, pero no se obtuvo ninguno")
	}
	if book != nil {
		t.Error("Se esperaba nil, pero se obtuvo un libro")
	}
	expectedError := "el autor del libro es obligatorio"
	if err.Error() != expectedError {
		t.Errorf("Se esperaba error '%s', pero se obtuvo: %s", expectedError, err.Error())
	}
}

// TestCreateBook_RepositoryError prueba el manejo de errores del repositorio
func TestCreateBook_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := NewMockBookRepository()
	mockRepo.SetShouldError(true) // Configurar el mock para que retorne error
	bookUseCase := usecase.NewBookUseCase(mockRepo)

	// Act
	book, err := bookUseCase.CreateBook("Título válido", "Autor válido")

	// Assert
	if err == nil {
		t.Error("Se esperaba un error del repositorio, pero no se obtuvo ninguno")
	}
	if book != nil {
		t.Error("Se esperaba nil, pero se obtuvo un libro")
	}
}

// TestGetBookByID_Success prueba obtener un libro exitosamente
func TestGetBookByID_Success(t *testing.T) {
	// Arrange
	mockRepo := NewMockBookRepository()
	bookUseCase := usecase.NewBookUseCase(mockRepo)

	// Primero crear un libro
	createdBook, _ := bookUseCase.CreateBook("Test Book", "Test Author")

	// Act
	foundBook, err := bookUseCase.GetBookByID(createdBook.ID)

	// Assert
	if err != nil {
		t.Errorf("Se esperaba que no hubiera error, pero se obtuvo: %v", err)
	}
	if foundBook == nil {
		t.Fatal("Se esperaba un libro, pero se obtuvo nil")
	}
	if foundBook.ID != createdBook.ID {
		t.Errorf("Se esperaba ID '%s', pero se obtuvo: %s", createdBook.ID, foundBook.ID)
	}
}

// TestGetBookByID_EmptyID prueba el error cuando el ID está vacío
func TestGetBookByID_EmptyID(t *testing.T) {
	// Arrange
	mockRepo := NewMockBookRepository()
	bookUseCase := usecase.NewBookUseCase(mockRepo)

	// Act
	book, err := bookUseCase.GetBookByID("")

	// Assert
	if err == nil {
		t.Error("Se esperaba un error, pero no se obtuvo ninguno")
	}
	if book != nil {
		t.Error("Se esperaba nil, pero se obtuvo un libro")
	}
	expectedError := "ID del libro es obligatorio"
	if err.Error() != expectedError {
		t.Errorf("Se esperaba error '%s', pero se obtuvo: %s", expectedError, err.Error())
	}
}

// TestGetAllBooks_Success prueba obtener todos los libros
func TestGetAllBooks_Success(t *testing.T) {
	// Arrange
	mockRepo := NewMockBookRepository()
	bookUseCase := usecase.NewBookUseCase(mockRepo)

	// Crear algunos libros de prueba
	bookUseCase.CreateBook("Libro 1", "Autor 1")
	bookUseCase.CreateBook("Libro 2", "Autor 2")

	// Act
	books, err := bookUseCase.GetAllBooks()

	// Assert
	if err != nil {
		t.Errorf("Se esperaba que no hubiera error, pero se obtuvo: %v", err)
	}
	if len(books) != 2 {
		t.Errorf("Se esperaban 2 libros, pero se obtuvieron: %d", len(books))
	}
}

// Para ejecutar estos tests, usa:
// go test ./internal/usecase/test -v
//
// Salida esperada:
// === RUN   TestCreateBook_Success
// --- PASS: TestCreateBook_Success (0.00s)
// === RUN   TestCreateBook_EmptyTitle
// --- PASS: TestCreateBook_EmptyTitle (0.00s)
// === RUN   TestCreateBook_EmptyAuthor
// --- PASS: TestCreateBook_EmptyAuthor (0.00s)
// === RUN   TestCreateBook_RepositoryError
// --- PASS: TestCreateBook_RepositoryError (0.00s)
// === RUN   TestGetBookByID_Success
// --- PASS: TestGetBookByID_Success (0.00s)
// === RUN   TestGetBookByID_EmptyID
// --- PASS: TestGetBookByID_EmptyID (0.00s)
// === RUN   TestGetAllBooks_Success
// --- PASS: TestGetAllBooks_Success (0.00s)
// PASS

// 💡 CONSEJOS PARA TESTING EN CLEAN ARCHITECTURE:
//
// 1. 🧪 Testea cada capa independientemente
// 2. 🎭 Usa mocks para dependencias externas
// 3. 📝 Tests descriptivos: TestFunction_Scenario_ExpectedBehavior
// 4. ✅ Cubre casos exitosos y de error
// 5. 🚀 Tests rápidos: sin I/O, sin sleep, sin dependencias externas
//
// 🌟 PATRONES DE TESTING ÚTILES:
// - Table-driven tests: probar múltiples casos con una tabla
// - Golden files: comparar outputs complejos con archivos de referencia
// - Test doubles: mocks, stubs, fakes, spies
// - Property-based testing: generar casos de prueba automáticamente
//
// 🎓 PREGUNTA PARA REFLEXIONAR:
// ¿Por qué estos tests son tan rápidos y confiables?
// Respuesta: Porque testean solo la lógica de negocio,
// sin dependencias externas como HTTP, base de datos, filesystem, etc.
