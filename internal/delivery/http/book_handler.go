package http

import (
	"go-book-clean-architecture-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

// BookHandler maneja las peticiones HTTP relacionadas con libros
// Los handlers son parte de la capa de delivery/infrastructure
// Su responsabilidad es convertir peticiones HTTP en llamadas a casos de uso
// y convertir las respuestas de los casos de uso en respuestas HTTP
type BookHandler struct {
	bookUseCase *usecase.BookUseCase // Dependencia inyectada del caso de uso
}

// NewBookHandler crea un nuevo handler para libros
func NewBookHandler(bookUseCase *usecase.BookUseCase) *BookHandler {
	return &BookHandler{
		bookUseCase: bookUseCase,
	}
}

// CreateBookRequest representa la estructura de datos esperada para crear un libro
type CreateBookRequest struct {
	Title  string `json:"title"`  // Título del libro
	Author string `json:"author"` // Autor del libro
}

// UpdateBookRequest representa la estructura de datos esperada para actualizar un libro
type UpdateBookRequest struct {
	Title  string `json:"title"`  // Título del libro
	Author string `json:"author"` // Autor del libro
}

// CreateBook maneja las peticiones POST /api/books
// Convierte la petición HTTP en una llamada al caso de uso
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	// Parsear el body de la petición
	var req CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de petición inválido",
		})
	}

	// Llamar al caso de uso
	book, err := h.bookUseCase.CreateBook(req.Title, req.Author)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retornar respuesta exitosa
	return c.Status(fiber.StatusCreated).JSON(book)
}

// GetBookByID maneja las peticiones GET /api/books/:id
func (h *BookHandler) GetBookByID(c *fiber.Ctx) error {
	// Obtener el ID del parámetro de la URL
	id := c.Params("id")

	// Llamar al caso de uso
	book, err := h.bookUseCase.GetBookByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retornar respuesta exitosa
	return c.JSON(book)
}

// GetAllBooks maneja las peticiones GET /api/books
func (h *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	// Llamar al caso de uso
	books, err := h.bookUseCase.GetAllBooks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retornar respuesta exitosa
	return c.JSON(books)
}

// UpdateBook maneja las peticiones PUT /api/books/:id
func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	// Obtener el ID del parámetro de la URL
	id := c.Params("id")

	// Parsear el body de la petición
	var req UpdateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de petición inválido",
		})
	}

	// Llamar al caso de uso
	book, err := h.bookUseCase.UpdateBook(id, req.Title, req.Author)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retornar respuesta exitosa
	return c.JSON(book)
}

// DeleteBook maneja las peticiones DELETE /api/books/:id
func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	// Obtener el ID del parámetro de la URL
	id := c.Params("id")

	// Llamar al caso de uso
	err := h.bookUseCase.DeleteBook(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retornar respuesta exitosa
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// UserHandler maneja las peticiones HTTP relacionadas con usuarios
type UserHandler struct {
	userUseCase *usecase.UserUseCase // Dependencia inyectada del caso de uso
}

// NewUserHandler crea un nuevo handler para usuarios
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
	// Obtener el ID del parámetro de la URL
	id := c.Params("id")

	// Llamar al caso de uso
	user, err := h.userUseCase.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retornar respuesta exitosa
	return c.JSON(user)
}

// GetAllUsers maneja las peticiones GET /api/users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Llamar al caso de uso
	users, err := h.userUseCase.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retornar respuesta exitosa
	return c.JSON(users)
}

// UpdateUser maneja las peticiones PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Obtener el ID del parámetro de la URL
	id := c.Params("id")

	// Parsear el body de la petición
	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de petición inválido",
		})
	}

	// Llamar al caso de uso
	user, err := h.userUseCase.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retornar respuesta exitosa
	return c.JSON(user)
}

// DeleteUser maneja las peticiones DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// Obtener el ID del parámetro de la URL
	id := c.Params("id")

	// Llamar al caso de uso
	err := h.userUseCase.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retornar respuesta exitosa
	return c.Status(fiber.StatusNoContent).Send(nil)
}
