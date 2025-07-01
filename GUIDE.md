# 📖 Guía Paso a Paso: Clean Architecture con Go y Fiber

## 🎯 Objetivo
Esta guía te enseña cómo crear una API REST siguiendo Clean Architecture, paso a paso, explicando cada decisión y patrón utilizado.

## 📚 ¿Qué vamos a aprender?

1. **Clean Architecture** - Organización del código en capas
2. **Dependency Injection** - Cómo inyectar dependencias
3. **Repository Pattern** - Abstracción de la persistencia
4. **Use Cases** - Lógica de negocio separada
5. **HTTP Handlers** - Manejo de peticiones web
6. **Fiber Framework** - Framework web para Go

## 🏗️ Las 4 Capas de Clean Architecture

```
┌─────────────────────────────────────┐
│         CAPA EXTERNA                │
│  (Delivery/Infrastructure)          │
│  ┌─────────────────────────────┐    │
│  │      CAPA INTERFAZ          │    │
│  │     (Use Cases)             │    │
│  │  ┌─────────────────────┐    │    │
│  │  │   CAPA DOMINIO      │    │    │
│  │  │   (Entities)        │    │    │
│  │  └─────────────────────┘    │    │
│  └─────────────────────────────┘    │
└─────────────────────────────────────┘
```

### 🎯 Regla de Dependencias
**Las dependencias apuntan hacia adentro**
- External → Interface → Application → Domain
- Nunca al revés

## 🚀 Paso 1: Configurar el Proyecto

### 1.1 Crear el módulo Go
```bash
mkdir go-book-clean-architecture-api
cd go-book-clean-architecture-api
go mod init go-book-clean-architecture-api
```

### 1.2 Instalar dependencias
```bash
go get github.com/gofiber/fiber/v2
go get github.com/google/uuid
```

### 1.3 Crear estructura de carpetas
```bash
mkdir -p cmd/server
mkdir -p internal/domain
mkdir -p internal/usecase
mkdir -p internal/repository
mkdir -p internal/delivery/http
mkdir -p internal/routes
mkdir -p internal/infrastructure/memory
```

## 📖 Paso 2: Crear las Entidades (Domain Layer)

**Archivo:** `internal/domain/book.go`

```go
package domain

// Book representa la entidad principal de nuestro dominio
// En Clean Architecture, las entidades contienen la lógica de negocio más crítica
// y son independientes de frameworks, bases de datos, UI, etc.
type Book struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

// User representa la entidad de usuario en nuestro dominio
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

**¿Por qué esta capa?**
- ✅ Contiene las reglas de negocio más importantes
- ✅ Es independiente de todo (frameworks, DB, UI)
- ✅ Es el corazón de la aplicación

## 📋 Paso 3: Definir Contratos (Repository Interfaces)

**Archivo:** `internal/repository/book_repository.go`

```go
package repository

import "go-book-clean-architecture-api/internal/domain"

// BookRepository define el contrato para las operaciones de persistencia
// Esta interfaz está en la capa de aplicación, pero las implementaciones
// están en la capa de infraestructura (Dependency Inversion Principle)
type BookRepository interface {
    Create(book *domain.Book) (*domain.Book, error)
    GetByID(id string) (*domain.Book, error)
    GetAll() ([]*domain.Book, error)
    Update(book *domain.Book) (*domain.Book, error)
    Delete(id string) error
}
```

**¿Por qué interfaces?**
- ✅ Permiten cambiar implementaciones sin romper código
- ✅ Facilitan el testing (mocks)
- ✅ Implementan el Dependency Inversion Principle

## 🧠 Paso 4: Implementar Casos de Uso (Application Layer)

**Archivo:** `internal/usecase/book_usecase.go`

```go
package usecase

import (
    "errors"
    "go-book-clean-architecture-api/internal/domain"
    "go-book-clean-architecture-api/internal/repository"
    "github.com/google/uuid"
)

// BookUseCase contiene la lógica de negocio de la aplicación
type BookUseCase struct {
    bookRepo repository.BookRepository // Dependencia inyectada
}

// NewBookUseCase constructor con inyección de dependencias
func NewBookUseCase(bookRepo repository.BookRepository) *BookUseCase {
    return &BookUseCase{
        bookRepo: bookRepo,
    }
}

// CreateBook implementa la lógica para crear un libro
func (uc *BookUseCase) CreateBook(title, author string) (*domain.Book, error) {
    // Validaciones de negocio
    if title == "" {
        return nil, errors.New("el título del libro es obligatorio")
    }
    if author == "" {
        return nil, errors.New("el autor del libro es obligatorio")
    }

    // Crear entidad
    book := &domain.Book{
        ID:     uuid.New().String(),
        Title:  title,
        Author: author,
    }

    // Delegar persistencia al repositorio
    return uc.bookRepo.Create(book)
}
```

**¿Por qué casos de uso?**
- ✅ Contienen la lógica específica de la aplicación
- ✅ Orquestan operaciones entre entidades
- ✅ Definen qué se puede hacer en la aplicación

## 💾 Paso 5: Implementar la Persistencia (Infrastructure Layer)

**Archivo:** `internal/infrastructure/memory/book_repository.go`

```go
package memory

import (
    "errors"
    "go-book-clean-architecture-api/internal/domain"
    "go-book-clean-architecture-api/internal/repository"
    "sync"
)

// InMemoryBookRepository implementa BookRepository en memoria
type InMemoryBookRepository struct {
    books map[string]*domain.Book
    mutex sync.RWMutex
}

// NewInMemoryBookRepository constructor
func NewInMemoryBookRepository() repository.BookRepository {
    return &InMemoryBookRepository{
        books: make(map[string]*domain.Book),
        mutex: sync.RWMutex{},
    }
}

// Create implementa la persistencia de un libro
func (r *InMemoryBookRepository) Create(book *domain.Book) (*domain.Book, error) {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.books[book.ID]; exists {
        return nil, errors.New("el libro ya existe")
    }

    r.books[book.ID] = book
    return book, nil
}
```

**¿Por qué esta implementación?**
- ✅ Implementa la interfaz definida en casos de uso
- ✅ Se puede cambiar por PostgreSQL, MongoDB, etc.
- ✅ Está aislada del resto de la aplicación

## 🌐 Paso 6: Crear Handlers HTTP (Delivery Layer)

**Archivo:** `internal/delivery/http/book_handler.go`

```go
package http

import (
    "go-book-clean-architecture-api/internal/usecase"
    "github.com/gofiber/fiber/v2"
)

// BookHandler maneja peticiones HTTP para libros
type BookHandler struct {
    bookUseCase *usecase.BookUseCase
}

// NewBookHandler constructor
func NewBookHandler(bookUseCase *usecase.BookUseCase) *BookHandler {
    return &BookHandler{
        bookUseCase: bookUseCase,
    }
}

// CreateBook maneja POST /api/books
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
    var req struct {
        Title  string `json:"title"`
        Author string `json:"author"`
    }

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Formato inválido",
        })
    }

    book, err := h.bookUseCase.CreateBook(req.Title, req.Author)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(book)
}
```

**¿Por qué handlers?**
- ✅ Convierten peticiones HTTP en llamadas a casos de uso
- ✅ Manejan serialización/deserialización
- ✅ Están acoplados al framework web (pero aislados)

## 🛣️ Paso 7: Configurar Rutas

**Archivo:** `internal/routes/book_routes.go`

```go
package routes

import (
    "go-book-clean-architecture-api/internal/delivery/http"
    "github.com/gofiber/fiber/v2"
)

// SetupRoutes configura todas las rutas
func SetupRoutes(app *fiber.App, bookHandler *http.BookHandler) {
    api := app.Group("/api")
    
    books := api.Group("/books")
    books.Post("/", bookHandler.CreateBook)
    books.Get("/", bookHandler.GetAllBooks)
    books.Get("/:id", bookHandler.GetBookByID)
    books.Put("/:id", bookHandler.UpdateBook)
    books.Delete("/:id", bookHandler.DeleteBook)
}
```

## 🚀 Paso 8: Punto de Entrada y Dependency Injection

**Archivo:** `cmd/server/main.go`

```go
package main

import (
    "log"
    "go-book-clean-architecture-api/internal/delivery/http"
    "go-book-clean-architecture-api/internal/infrastructure/memory"
    "go-book-clean-architecture-api/internal/routes"
    "go-book-clean-architecture-api/internal/usecase"
    "github.com/gofiber/fiber/v2"
)

func main() {
    // 1. Crear aplicación Fiber
    app := fiber.New()

    // 2. DEPENDENCY INJECTION
    // 2.1. Crear repositorios (Infraestructura)
    bookRepo := memory.NewInMemoryBookRepository()
    
    // 2.2. Crear casos de uso (Aplicación)
    bookUseCase := usecase.NewBookUseCase(bookRepo)
    
    // 2.3. Crear handlers (Delivery)
    bookHandler := http.NewBookHandler(bookUseCase)

    // 3. Configurar rutas
    routes.SetupRoutes(app, bookHandler)

    // 4. Iniciar servidor
    log.Fatal(app.Listen(":8080"))
}
```

## 🔄 Flujo de una Petición

```
1. HTTP Request → 2. Handler → 3. Use Case → 4. Repository → 5. Database
                ↓
6. HTTP Response ← 5. Handler ← 4. Use Case ← 3. Repository ← 2. Database
```

### Ejemplo: Crear un libro

1. **Cliente:** `POST /api/books {"title": "Go Book", "author": "John"}`
2. **Handler:** Parsea JSON → llama `bookUseCase.CreateBook()`
3. **Use Case:** Valida datos → crea entidad → llama `bookRepo.Create()`
4. **Repository:** Guarda en memoria/DB
5. **Respuesta:** Se devuelve el libro creado

## 🧪 Testing con Clean Architecture

### Test del Use Case (sin HTTP, sin DB)
```go
func TestCreateBook(t *testing.T) {
    // Arrange
    mockRepo := &MockBookRepository{}
    useCase := usecase.NewBookUseCase(mockRepo)
    
    // Act
    book, err := useCase.CreateBook("Title", "Author")
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "Title", book.Title)
}
```

### Test del Handler (mock use case)
```go
func TestCreateBookHandler(t *testing.T) {
    // Arrange
    mockUseCase := &MockBookUseCase{}
    handler := http.NewBookHandler(mockUseCase)
    
    // Act & Assert usando Fiber test utilities
}
```

## 🎯 Beneficios de esta Arquitectura

### ✅ Testeable
- Cada capa se puede testear independientemente
- Fácil crear mocks para las interfaces

### ✅ Mantenible
- Cambios en una capa no afectan otras
- Código organizado y predecible

### ✅ Escalable
- Fácil agregar nuevas funcionalidades
- Separación clara de responsabilidades

### ✅ Flexible
- Cambiar de Fiber a Gin: solo cambiar handlers
- Cambiar de memoria a PostgreSQL: solo cambiar repositorio
- Agregar GraphQL: agregar nueva capa de delivery

## 🤔 Preguntas Frecuentes

**P: ¿Es mucho código para algo simple?**
R: Al inicio sí, pero cuando el proyecto crezca, agradecerás la organización.

**P: ¿Cuándo usar Clean Architecture?**
R: En proyectos que van a crecer, tienen múltiples desarrolladores, o necesitan ser mantenidos a largo plazo.

**P: ¿Es obligatorio seguir esta estructura exacta?**
R: No, puedes adaptarla. Lo importante son los principios: separación de capas y inversión de dependencias.

## 🎓 Próximos Pasos

1. **Agrega tests** a cada capa
2. **Implementa una base de datos real** (PostgreSQL)
3. **Agrega middleware** (autenticación, logging)
4. **Implementa validaciones** más robustas
5. **Agrega documentación** con Swagger

## 📚 Recursos Adicionales

- [Clean Architecture - Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Dependency Injection in Go](https://blog.drewolson.org/dependency-injection-in-go)
- [Testing in Go](https://go.dev/doc/tutorial/add-a-test)

---

**¡Feliz aprendizaje! 🚀**
