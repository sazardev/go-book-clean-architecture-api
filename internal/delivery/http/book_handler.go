// Package http contiene los HANDLERS que manejan peticiones HTTP
// Esta es la CAPA DE DELIVERY/INTERFAZ de Clean Architecture
//
// 🌐 ¿Qué son los handlers HTTP?
// - Son adaptadores que convierten peticiones HTTP en llamadas a casos de uso
// - Manejan la serialización/deserialización de datos (JSON ↔ Go structs)
// - Se encargan de los códigos de estado HTTP apropiados
// - Son específicos del protocolo HTTP (a diferencia de los casos de uso)
//
// 🎯 REGLAS DE ORO de los Handlers:
// ✅ Solo manejan conversión entre HTTP y casos de uso
// ✅ NO contienen lógica de negocio (esa va en casos de uso)
// ✅ Validan formato de entrada, pero NO reglas de negocio
// ✅ Convierten errores de casos de uso en códigos HTTP apropiados
// ✅ Son específicos del framework web (en este caso, Fiber)
//
// 🚫 Los handlers NO deben:
// - Contener validaciones de negocio
// - Acceder directamente a repositorios
// - Conocer detalles de base de datos
// - Tener lógica compleja de transformación de datos
package http

import (
	"go-book-clean-architecture-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

// BookHandler maneja las peticiones HTTP relacionadas con libros
//
// 📚 ¿Por qué separamos BookHandler de UserHandler?
// - Principio de Responsabilidad Única (SRP)
// - Cada handler se enfoca en una entidad específica
// - Más fácil de mantener y testear
// - Permite escalabilidad (diferentes equipos pueden trabajar en diferentes handlers)
//
// 🔧 Patrón utilizado: Dependency Injection
// - El handler RECIBE el caso de uso que necesita
// - NO crea el caso de uso internamente
// - Esto facilita el testing (podemos inyectar mocks)
type BookHandler struct {
	bookUseCase *usecase.BookUseCase // Dependencia inyectada del caso de uso
}

// NewBookHandler es el CONSTRUCTOR que implementa Dependency Injection
//
// 🎯 ¿Por qué usar constructores para handlers?
// - Garantizan que el handler se crea con todas las dependencias necesarias
// - Hacen explícitas las dependencias (fácil ver qué necesita cada handler)
// - Facilitan el testing
// - Siguen las mejores prácticas de Go
func NewBookHandler(bookUseCase *usecase.BookUseCase) *BookHandler {
	return &BookHandler{
		bookUseCase: bookUseCase,
	}
}

// CreateBookRequest representa la estructura de datos esperada para crear un libro
//
// 📝 ¿Por qué definir structs para requests?
// - Tipado fuerte: Go puede validar la estructura automáticamente
// - Documentación: es claro qué campos se esperan
// - Mantenimiento: fácil agregar/quitar campos
// - Reutilización: el mismo struct se puede usar en testing
//
// 🏷️ Tags JSON: definen cómo se serializa/deserializa desde/hacia JSON
type CreateBookRequest struct {
	Title  string `json:"title"`  // Título del libro
	Author string `json:"author"` // Autor del libro
}

// UpdateBookRequest representa la estructura de datos esperada para actualizar un libro
// Nota: Mismo contenido que CreateBookRequest, pero semánticamente diferente
type UpdateBookRequest struct {
	Title  string `json:"title"`  // Título del libro
	Author string `json:"author"` // Autor del libro
}

// CreateBook maneja las peticiones POST /api/books
//
// 🔄 Flujo típico de un handler HTTP:
// 1. Parsear y validar la petición HTTP
// 2. Extraer datos necesarios
// 3. Llamar al caso de uso correspondiente
// 4. Convertir la respuesta del caso de uso a HTTP
// 5. Retornar respuesta con código de estado apropiado
//
// 📊 Códigos de estado HTTP utilizados:
// - 201 Created: recurso creado exitosamente
// - 400 Bad Request: formato de petición inválido o error de validación
// - 500 Internal Server Error: error interno del servidor
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	// PASO 1: Parsear el body de la petición HTTP
	var req CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		// Error de formato: el JSON no es válido o no coincide con el struct
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de petición inválido",
		})
	}

	// PASO 2: Llamar al caso de uso (aquí es donde ocurre la magia)
	// El handler NO valida reglas de negocio, solo delega al caso de uso
	book, err := h.bookUseCase.CreateBook(req.Title, req.Author)
	if err != nil {
		// Error de negocio: título vacío, autor vacío, etc.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// PASO 3: Retornar respuesta exitosa
	// 201 Created es el código apropiado para creación de recursos
	return c.Status(fiber.StatusCreated).JSON(book)
}

// GetBookByID maneja las peticiones GET /api/books/:id
//
// 🔍 Handler para obtener un recurso específico
// Utiliza parámetros de URL para obtener el ID
func (h *BookHandler) GetBookByID(c *fiber.Ctx) error {
	// PASO 1: Obtener el ID del parámetro de la URL
	// :id en la ruta se convierte en un parámetro accesible
	id := c.Params("id")

	// PASO 2: Llamar al caso de uso
	book, err := h.bookUseCase.GetBookByID(id)
	if err != nil {
		// 404 Not Found es apropiado cuando el recurso no existe
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// PASO 3: Retornar respuesta exitosa
	// 200 OK es el código por defecto para consultas exitosas
	return c.JSON(book)
}

// GetAllBooks maneja las peticiones GET /api/books
//
// 📚 Handler para obtener una colección de recursos
// En aplicaciones reales, implementarías paginación aquí
func (h *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	// PASO 1: Llamar al caso de uso
	// No necesitamos parámetros para obtener todos los libros
	books, err := h.bookUseCase.GetAllBooks()
	if err != nil {
		// 500 Internal Server Error para errores inesperados
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// PASO 2: Retornar respuesta exitosa
	// Nota: si no hay libros, retornamos un array vacío, no un error
	return c.JSON(books)
}

// UpdateBook maneja las peticiones PUT /api/books/:id
//
// ✏️ Handler para actualizar un recurso existente
// Combina parámetros de URL (ID) con body de petición (datos)
func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	// PASO 1: Obtener el ID del parámetro de la URL
	id := c.Params("id")

	// PASO 2: Parsear el body de la petición
	var req UpdateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de petición inválido",
		})
	}

	// PASO 3: Llamar al caso de uso
	book, err := h.bookUseCase.UpdateBook(id, req.Title, req.Author)
	if err != nil {
		// Podría ser 400 (validación) o 404 (no existe)
		// En este caso, simplificamos con 400
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// PASO 4: Retornar respuesta exitosa
	// 200 OK es apropiado para actualizaciones exitosas
	return c.JSON(book)
}

// DeleteBook maneja las peticiones DELETE /api/books/:id
//
// 🗑️ Handler para eliminar un recurso
// Retorna 204 No Content en caso de éxito
func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	// PASO 1: Obtener el ID del parámetro de la URL
	id := c.Params("id")

	// PASO 2: Llamar al caso de uso
	err := h.bookUseCase.DeleteBook(id)
	if err != nil {
		// 404 Not Found si el libro no existe
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// PASO 3: Retornar respuesta exitosa sin contenido
	// 204 No Content es el código estándar para eliminación exitosa
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// UserHandler maneja las peticiones HTTP relacionadas con usuarios
//
// 👤 Misma estructura que BookHandler, pero para usuarios
// Esto demuestra el patrón consistente en Clean Architecture
type UserHandler struct {
	userUseCase *usecase.UserUseCase // Dependencia inyectada del caso de uso
}

// NewUserHandler constructor para UserHandler
func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUserRequest representa la estructura de datos esperada para crear un usuario
type CreateUserRequest struct {
	Name  string `json:"name"`  // Nombre del usuario
	Email string `json:"email"` // Email del usuario
}

// UpdateUserRequest representa la estructura de datos esperada para actualizar un usuario
type UpdateUserRequest struct {
	Name  string `json:"name"`  // Nombre del usuario
	Email string `json:"email"` // Email del usuario
}

// CreateUser maneja las peticiones POST /api/users
//
// 👤 Mismo patrón que CreateBook, pero para usuarios
// La consistencia en los patrones facilita el mantenimiento
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Parsear el body de la petición
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de petición inválido",
		})
	}

	// Llamar al caso de uso
	user, err := h.userUseCase.CreateUser(req.Name, req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retornar respuesta exitosa
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUserByID maneja las peticiones GET /api/users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userUseCase.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user)
}

// GetAllUsers maneja las peticiones GET /api/users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userUseCase.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(users)
}

// UpdateUser maneja las peticiones PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de petición inválido",
		})
	}

	user, err := h.userUseCase.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user)
}

// DeleteUser maneja las peticiones DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.userUseCase.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// 💡 CONSEJOS PARA PRINCIPIANTES:
//
// 1. 🎯 Un handler = Un endpoint específico
// 2. 📝 Siempre validar formato de entrada (BodyParser)
// 3. 🔧 Delegar toda la lógica de negocio a los casos de uso
// 4. 📊 Usar códigos de estado HTTP apropiados
// 5. 🚫 NO poner lógica de negocio en los handlers
//
// 🌟 CÓDIGOS DE ESTADO HTTP MÁS COMUNES:
// - 200 OK: operación exitosa
// - 201 Created: recurso creado exitosamente
// - 204 No Content: operación exitosa sin contenido de respuesta
// - 400 Bad Request: error en la petición del cliente
// - 404 Not Found: recurso no encontrado
// - 500 Internal Server Error: error interno del servidor
//
// 🚫 EJEMPLOS DE LO QUE NO DEBES PONER AQUÍ:
// - Conexiones directas a base de datos
// - Validaciones de reglas de negocio complejas
// - Lógica de cálculos o transformaciones de datos
// - Dependencias de librerías específicas de persistencia
//
// 🎓 PREGUNTA PARA REFLEXIONAR:
// Si quisiéramos cambiar de Fiber a Gin (otro framework web),
// ¿qué tendríamos que cambiar?
// Respuesta: Solo estos handlers, los casos de uso permanecerían iguales.
