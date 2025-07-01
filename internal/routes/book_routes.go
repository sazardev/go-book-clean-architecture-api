package routes

import (
	"go-book-clean-architecture-api/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

// SetupBookRoutes configura todas las rutas relacionadas con libros
// Las rutas definen qué handler se ejecuta para cada endpoint
// Esto separa la configuración de rutas de la lógica de los handlers
func SetupBookRoutes(app *fiber.App, bookHandler *http.BookHandler) {
	// Crear un grupo de rutas para libros con prefijo /api/books
	books := app.Group("/api/books")

	// Configurar las rutas CRUD para libros
	books.Post("/", bookHandler.CreateBook)      // POST /api/books - Crear libro
	books.Get("/", bookHandler.GetAllBooks)      // GET /api/books - Obtener todos los libros
	books.Get("/:id", bookHandler.GetBookByID)   // GET /api/books/:id - Obtener libro por ID
	books.Put("/:id", bookHandler.UpdateBook)    // PUT /api/books/:id - Actualizar libro
	books.Delete("/:id", bookHandler.DeleteBook) // DELETE /api/books/:id - Eliminar libro
}

// SetupUserRoutes configura todas las rutas relacionadas con usuarios
func SetupUserRoutes(app *fiber.App, userHandler *http.UserHandler) {
	// Crear un grupo de rutas para usuarios con prefijo /api/users
	users := app.Group("/api/users")

	// Configurar las rutas CRUD para usuarios
	users.Post("/", userHandler.CreateUser)      // POST /api/users - Crear usuario
	users.Get("/", userHandler.GetAllUsers)      // GET /api/users - Obtener todos los usuarios
	users.Get("/:id", userHandler.GetUserByID)   // GET /api/users/:id - Obtener usuario por ID
	users.Put("/:id", userHandler.UpdateUser)    // PUT /api/users/:id - Actualizar usuario
	users.Delete("/:id", userHandler.DeleteUser) // DELETE /api/users/:id - Eliminar usuario
}

// SetupRoutes configura todas las rutas de la aplicación
// Esta función central configura todos los endpoints de la API
func SetupRoutes(app *fiber.App, bookHandler *http.BookHandler, userHandler *http.UserHandler) {
	// Ruta de health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "API funcionando correctamente",
		})
	})

	// Configurar rutas específicas para cada dominio
	SetupBookRoutes(app, bookHandler)
	SetupUserRoutes(app, userHandler)
}
