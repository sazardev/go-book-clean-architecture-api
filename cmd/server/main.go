package main

import (
	"log"

	"go-book-clean-architecture-api/internal/delivery/http"
	"go-book-clean-architecture-api/internal/infrastructure/memory"
	"go-book-clean-architecture-api/internal/routes"
	"go-book-clean-architecture-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// PASO 1: Crear la aplicación Fiber
	// Fiber es nuestro framework web, pero está aislado en la capa de delivery
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// PASO 2: Configurar middleware básico
	app.Use(logger.New()) // Logging de peticiones HTTP
	app.Use(cors.New())   // Habilitar CORS para peticiones desde el frontend

	// PASO 3: INYECCIÓN DE DEPENDENCIAS - Aquí es donde ensamblamos toda la aplicación
	// Esta es la parte más importante de Clean Architecture

	// 3.1: Crear las implementaciones de los repositorios (capa de infraestructura)
	bookRepo := memory.NewInMemoryBookRepository() // Implementación en memoria
	userRepo := memory.NewInMemoryUserRepository() // Implementación en memoria

	// 3.2: Crear los casos de uso e inyectar los repositorios (capa de dominio/aplicación)
	bookUseCase := usecase.NewBookUseCase(bookRepo) // Inyectar repositorio de libros
	userUseCase := usecase.NewUserUseCase(userRepo) // Inyectar repositorio de usuarios

	// 3.3: Crear los handlers e inyectar los casos de uso (capa de delivery)
	bookHandler := http.NewBookHandler(bookUseCase) // Inyectar caso de uso de libros
	userHandler := http.NewUserHandler(userUseCase) // Inyectar caso de uso de usuarios

	// PASO 4: Configurar las rutas
	routes.SetupRoutes(app, bookHandler, userHandler)

	// PASO 5: Iniciar el servidor
	log.Println("🚀 Servidor iniciado en http://localhost:8080")
	log.Println("📚 Endpoints disponibles:")
	log.Println("  GET    /health           - Health check")
	log.Println("  POST   /api/books        - Crear libro")
	log.Println("  GET    /api/books        - Obtener todos los libros")
	log.Println("  GET    /api/books/:id    - Obtener libro por ID")
	log.Println("  PUT    /api/books/:id    - Actualizar libro")
	log.Println("  DELETE /api/books/:id    - Eliminar libro")
	log.Println("  POST   /api/users        - Crear usuario")
	log.Println("  GET    /api/users        - Obtener todos los usuarios")
	log.Println("  GET    /api/users/:id    - Obtener usuario por ID")
	log.Println("  PUT    /api/users/:id    - Actualizar usuario")
	log.Println("  DELETE /api/users/:id    - Eliminar usuario")

	// Iniciar el servidor en el puerto 8080
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}

/*
EXPLICACIÓN DEL FLUJO DE CLEAN ARCHITECTURE:

1. **Capa de Dominio (Domain)**:
   - Contiene las entidades principales (Book, User)
   - Define las reglas de negocio más importantes
   - Es independiente de frameworks y librerías externas

2. **Capa de Aplicación/Casos de Uso (Use Cases)**:
   - Contiene la lógica de negocio específica de la aplicación
   - Define las interfaces de los repositorios (contratos)
   - Orquesta las operaciones entre entidades

3. **Capa de Infraestructura (Infrastructure)**:
   - Implementa las interfaces definidas en casos de uso
   - Maneja la persistencia de datos (en este caso, en memoria)
   - Puede incluir conexiones a bases de datos, APIs externas, etc.

4. **Capa de Delivery/Interfaz (Delivery)**:
   - Maneja la comunicación con el exterior (HTTP, CLI, gRPC, etc.)
   - Convierte peticiones externas en llamadas a casos de uso
   - Convierte respuestas de casos de uso en formatos externos

5. **Inyección de Dependencias**:
   - Las dependencias fluyen hacia adentro (hacia el dominio)
   - Las capas externas dependen de las internas, nunca al revés
   - Esto permite cambiar implementaciones sin afectar la lógica de negocio

BENEFICIOS:
- ✅ Testeable: Cada capa se puede testear independientemente
- ✅ Mantenible: Cambios en una capa no afectan otras
- ✅ Escalable: Fácil agregar nuevas funcionalidades
- ✅ Flexible: Fácil cambiar frameworks, bases de datos, etc.
- ✅ Independiente: La lógica de negocio no depende de librerías externas
*/
