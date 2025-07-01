# 📚 Go Clean Architecture API - Guía Completa para Principiantes

> **¡La plantilla más simple y educativa de Clean Architecture en Go con Fiber!**

Este proyecto es una API RESTful súper simple para gestionar **libros** y **usuarios**, implementada siguiendo los principios de **Clean Architecture**. Está diseñado específicamente para **principiantes en Go** que quieren entender Clean Architecture de manera práctica y sencilla.

## 🎯 ¿Qué es Clean Architecture?

**Clean Architecture** es una forma de organizar tu código que:
- ✅ Separa la lógica de negocio del framework web
- ✅ Hace el código fácil de testear
- ✅ Permite cambiar la base de datos sin romper nada
- ✅ Mantiene el código organizado y escalable

## 🏗️ Estructura del Proyecto (Súper Simple)

```
go-book-clean-architecture-api/
│
├── cmd/server/main.go                    # 🚀 Punto de entrada (aquí arranca todo)
│
├── internal/
│   ├── domain/book.go                    # 📖 Entidades (Book y User)
│   ├── repository/book_repository.go     # 📋 Contratos (interfaces)
│   ├── usecase/book_usecause.go         # 🧠 Lógica de negocio
│   ├── delivery/http/book_handler.go     # 🌐 Handlers HTTP
│   ├── routes/book_routes.go             # 🛣️ Rutas de la API
│   └── infrastructure/memory/            # 💾 Implementación en memoria
│       └── book_repository.go
│
├── go.mod                                # 📦 Dependencias
├── README.md                             # 📖 Esta guía
├── api_examples.http                     # 🧪 Ejemplos de peticiones
└── GUIDE.md                              # 📝 Guía paso a paso
```

## 🚀 Cómo usar este proyecto

### 1. Clonar e instalar dependencias
```bash
git clone <este-repo>
cd go-book-clean-architecture-api
go mod tidy
```

### 2. Ejecutar la aplicación
```bash
go run cmd/server/main.go
```

### 3. Probar la API
El servidor arranca en `http://localhost:8080`

**Endpoints disponibles:**
- `GET /health` - Verificar que la API funciona
- `POST /api/books` - Crear un libro
- `GET /api/books` - Obtener todos los libros
- `GET /api/books/:id` - Obtener un libro específico
- `PUT /api/books/:id` - Actualizar un libro
- `DELETE /api/books/:id` - Eliminar un libro
- `POST /api/users` - Crear un usuario
- `GET /api/users` - Obtener todos los usuarios
- (Y más endpoints para usuarios...)

## 🧪 Ejemplos de uso

### Crear un libro
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{"title": "Clean Architecture", "author": "Robert C. Martin"}'
```

### Obtener todos los libros
```bash
curl http://localhost:8080/api/books
```

## 🎓 Guía de Aprendizaje (Las 4 Capas)

### 🏛️ 1. Capa de Dominio (`internal/domain/`)
**¿Qué es?** Las entidades principales de tu negocio.
**Archivo:** `book.go`
```go
type Book struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}
```
**Regla de oro:** Esta capa NO conoce nada de HTTP, bases de datos, o frameworks.

### 🧠 2. Capa de Casos de Uso (`internal/usecase/`)
**¿Qué es?** La lógica de negocio de tu aplicación.
**Archivo:** `book_usecause.go`
```go
func (uc *BookUseCase) CreateBook(title, author string) (*domain.Book, error) {
    // Validaciones de negocio aquí
    if title == "" {
        return nil, errors.New("el título es obligatorio")
    }
    // Crear y guardar el libro
}
```
**Regla de oro:** Aquí van las validaciones y la lógica de tu negocio.

### 💾 3. Capa de Infraestructura (`internal/infrastructure/`)
**¿Qué es?** Cómo guardas los datos (base de datos, archivos, memoria, etc.).
**Archivo:** `memory/book_repository.go`
```go
func (r *InMemoryBookRepository) Create(book *domain.Book) (*domain.Book, error) {
    // Guardar en memoria (en producción sería PostgreSQL, MongoDB, etc.)
}
```
**Regla de oro:** Implementa las interfaces definidas en los casos de uso.

### 🌐 4. Capa de Delivery (`internal/delivery/http/`)
**¿Qué es?** Cómo el mundo exterior interactúa con tu aplicación (HTTP, CLI, etc.).
**Archivo:** `book_handler.go`
```go
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
    // Convertir petición HTTP en llamada al caso de uso
}
```
**Regla de oro:** Solo maneja la conversión entre HTTP y casos de uso.

## 🔧 ¿Cómo agregar un nuevo endpoint?

### Ejemplo: Agregar endpoint para buscar libros por autor

**Paso 1:** Agregar método al repositorio (`repository/book_repository.go`)
```go
type BookRepository interface {
    // ...métodos existentes...
    GetByAuthor(author string) ([]*domain.Book, error)
}
```

**Paso 2:** Implementar en infraestructura (`infrastructure/memory/book_repository.go`)
```go
func (r *InMemoryBookRepository) GetByAuthor(author string) ([]*domain.Book, error) {
    // Implementar búsqueda por autor
}
```

**Paso 3:** Agregar caso de uso (`usecase/book_usecause.go`)
```go
func (uc *BookUseCase) GetBooksByAuthor(author string) ([]*domain.Book, error) {
    if author == "" {
        return nil, errors.New("autor es obligatorio")
    }
    return uc.bookRepo.GetByAuthor(author)
}
```

**Paso 4:** Agregar handler (`delivery/http/book_handler.go`)
```go
func (h *BookHandler) GetBooksByAuthor(c *fiber.Ctx) error {
    author := c.Query("author")
    books, err := h.bookUseCase.GetBooksByAuthor(author)
    // Manejar respuesta...
}
```

**Paso 5:** Agregar ruta (`routes/book_routes.go`)
```go
books.Get("/search", bookHandler.GetBooksByAuthor) // GET /api/books/search?author=...
```

## 🎯 Ventajas de esta arquitectura

1. **🧪 Fácil de testear:** Cada capa se puede testear independientemente
2. **🔄 Fácil de cambiar:** Puedes cambiar de Fiber a Gin sin tocar la lógica de negocio
3. **💾 Flexible:** Cambiar de memoria a PostgreSQL es cambiar una línea en `main.go`
4. **📈 Escalable:** Agregar nuevas funcionalidades es simple y predecible
5. **👥 Trabajo en equipo:** Diferentes desarrolladores pueden trabajar en diferentes capas

## 🤔 Preguntas frecuentes

**P: ¿Por qué tantas carpetas para algo tan simple?**
R: Al principio parece complejo, pero cuando tu proyecto crezca, te agradecerás tener todo organizado.

**P: ¿Es obligatorio usar esta estructura?**
R: No, pero es una práctica recomendada por la comunidad de Go para proyectos serios.

**P: ¿Puedo usar otra base de datos?**
R: ¡Sí! Solo crea una nueva implementación en `infrastructure/` y cámbiala en `main.go`.

**P: ¿Puedo usar otro framework web?**
R: ¡Sí! Solo cambia los handlers y mantén la misma estructura.

## 📚 Recursos adicionales

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go by Example](https://gobyexample.com/)
- [Fiber Documentation](https://docs.gofiber.io/)

## 🤝 Contribuir

¿Encontraste algo confuso? ¿Tienes sugerencias? ¡Abre un issue o pull request!

---

**¡Feliz coding! 🚀**