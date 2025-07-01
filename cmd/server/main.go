// Package main es el punto de entrada de la aplicación
// Este archivo es donde se ensambla TODA la aplicación usando Clean Architecture
//
// 🎯 RESPONSABILIDADES del main.go:
// ✅ Configurar el framework web (Fiber)
// ✅ Implementar DEPENDENCY INJECTION (¡La magia de Clean Architecture!)
// ✅ Conectar todas las capas: Domain → UseCase → Infrastructure → Delivery
// ✅ Configurar middleware básico
// ✅ Iniciar el servidor
//
// 🚫 El main.go NO debe:
// - Contener lógica de negocio
// - Tener validaciones complejas
// - Manejar datos específicos
//
// 💡 CONCEPTO CLAVE: Composition Root
// Este archivo es el "Composition Root" donde se ensambla toda la aplicación
// Aquí es donde las dependencias se "inyectan" siguiendo el patrón DI
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
	// 🎯 PASO 1: Crear la aplicación Fiber
	// Fiber es nuestro framework web, pero está completamente aislado en la capa de delivery
	// Si quisiéramos cambiar a Gin, Echo, etc., solo cambiaríamos esta línea y los handlers
	app := fiber.New(fiber.Config{
		// Configurar manejo global de errores
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			log.Printf("Error no manejado: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error interno del servidor",
			})
		},
		// Prefork para mejor performance en producción (opcional)
		Prefork: false,
		// Configuración de JSON más legible
		JSONEncoder: nil,
	})

	// 🎯 PASO 2: Configurar middleware básico
	// El middleware se ejecuta antes de llegar a los handlers
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	})) // Logging de todas las peticiones HTTP

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // En producción, especificar dominios exactos
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	})) // Habilitar CORS para peticiones desde el frontend

	// 🎯 PASO 3: DEPENDENCY INJECTION - ¡La parte MÁS IMPORTANTE!
	// Esta es la implementación práctica de Clean Architecture
	//
	// 🔄 FLUJO DE DEPENDENCIAS (de afuera hacia adentro):
	// Infrastructure → UseCase → Handler → Routes → App
	//
	// 📚 ¿Por qué este orden?
	// - Las capas internas NO conocen las externas
	// - Las capas externas SÍ conocen las internas
	// - Esto permite flexibilidad y testing

	log.Println("🔧 Configurando inyección de dependencias...")

	// 3.1: CAPA DE INFRAESTRUCTURA (más externa)
	// Aquí creamos las implementaciones concretas de persistencia
	log.Println("📁 Creando repositorios de infraestructura...")
	bookRepo := memory.NewInMemoryBookRepository() // Implementación en memoria
	userRepo := memory.NewInMemoryUserRepository() // Implementación en memoria

	// 💡 FLEXIBILIDAD: Para cambiar a PostgreSQL, solo cambiarías estas líneas:
	// bookRepo := postgresql.NewPostgresBookRepository(db)
	// userRepo := postgresql.NewPostgresUserRepository(db)

	log.Println("✅ Repositorios creados exitosamente")

	// 3.2: CAPA DE APLICACIÓN/CASOS DE USO (capa media)
	// Inyectamos los repositorios en los casos de uso
	log.Println("🧠 Creando casos de uso de aplicación...")
	bookUseCase := usecase.NewBookUseCase(bookRepo) // Inyectar repositorio de libros
	userUseCase := usecase.NewUserUseCase(userRepo) // Inyectar repositorio de usuarios

	log.Println("✅ Casos de uso creados exitosamente")

	// 3.3: CAPA DE DELIVERY/INTERFAZ (más interna de las externas)
	// Inyectamos los casos de uso en los handlers
	log.Println("🌐 Creando handlers de delivery...")
	bookHandler := http.NewBookHandler(bookUseCase) // Inyectar caso de uso de libros
	userHandler := http.NewUserHandler(userUseCase) // Inyectar caso de uso de usuarios

	log.Println("✅ Handlers creados exitosamente")

	// 🎯 PASO 4: Configurar las rutas
	// Las rutas conectan URLs con handlers específicos
	log.Println("🛣️ Configurando rutas de la aplicación...")
	routes.SetupRoutes(app, bookHandler, userHandler)
	log.Println("✅ Rutas configuradas exitosamente")

	// 🎯 PASO 5: Mostrar información útil y iniciar el servidor
	log.Println("")
	log.Println("🚀 ===== SERVIDOR INICIADO EXITOSAMENTE =====")
	log.Println("🌐 URL: http://localhost:8080")
	log.Println("� Documentación: README.md")
	log.Println("🧪 Ejemplos de peticiones: api_examples.http")
	log.Println("")
	log.Println("📚 ===== ENDPOINTS DISPONIBLES =====")
	log.Println("🔍 Health Check:")
	log.Println("  GET    /health              - Verificar estado de la API")
	log.Println("")
	log.Println("📖 Gestión de Libros:")
	log.Println("  POST   /api/books           - Crear un nuevo libro")
	log.Println("  GET    /api/books           - Obtener todos los libros")
	log.Println("  GET    /api/books/:id       - Obtener libro por ID")
	log.Println("  PUT    /api/books/:id       - Actualizar libro existente")
	log.Println("  DELETE /api/books/:id       - Eliminar libro")
	log.Println("")
	log.Println("👤 Gestión de Usuarios:")
	log.Println("  POST   /api/users           - Crear un nuevo usuario")
	log.Println("  GET    /api/users           - Obtener todos los usuarios")
	log.Println("  GET    /api/users/:id       - Obtener usuario por ID")
	log.Println("  PUT    /api/users/:id       - Actualizar usuario existente")
	log.Println("  DELETE /api/users/:id       - Eliminar usuario")
	log.Println("")
	log.Println("🎯 ===== EMPEZAR A PROBAR =====")
	log.Println("1. Abre api_examples.http en VS Code")
	log.Println("2. Instala la extensión 'REST Client'")
	log.Println("3. Haz clic en 'Send Request' en cualquier ejemplo")
	log.Println("4. ¡Experimenta y aprende!")
	log.Println("")

	// Iniciar el servidor en el puerto 8080
	// Esta línea bloquea el programa hasta que el servidor se detenga
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("💥 Error al iniciar el servidor:", err)
	}
}

/*
🎓 EXPLICACIÓN DETALLADA DEL FLUJO DE CLEAN ARCHITECTURE:

┌─────────────────────────────────────────────────────────────────┐
│                        CLEAN ARCHITECTURE                        │
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   DELIVERY  │────│ APPLICATION │────│   DOMAIN    │         │
│  │   (HTTP)    │    │ (USE CASES) │    │ (ENTITIES)  │         │
│  └─────────────┘    └─────────────┘    └─────────────┘         │
│         │                   │                   │              │
│         └───────────────────┼───────────────────┘              │
│                             │                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              INFRASTRUCTURE                             │   │
│  │            (REPOSITORIES)                               │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘

🔄 FLUJO DE UNA PETICIÓN:

1. 🌐 Cliente HTTP → 📨 POST /api/books
2. 🛣️ Router → 📍 Encuentra la ruta correspondiente
3. 🔗 Handler → 📝 bookHandler.CreateBook()
4. 🧠 Use Case → 💼 bookUseCase.CreateBook()
5. 💾 Repository → 🗃️ bookRepo.Create()
6. 📦 Response ← 🔙 Se devuelve la respuesta

🎯 BENEFICIOS DE ESTA ARQUITECTURA:

✅ TESTEABLE:
   - Cada capa se puede testear independientemente
   - Fácil crear mocks para las interfaces
   - Testing de casos de uso sin HTTP ni DB

✅ MANTENIBLE:
   - Cambios en una capa no afectan otras
   - Código organizado y predecible
   - Fácil localizar y corregir bugs

✅ ESCALABLE:
   - Fácil agregar nuevas funcionalidades
   - Separación clara de responsabilidades
   - Múltiples desarrolladores pueden trabajar en paralelo

✅ FLEXIBLE:
   - Cambiar de Fiber a Gin: solo cambiar handlers
   - Cambiar de memoria a PostgreSQL: solo cambiar repositorios
   - Agregar GraphQL: agregar nueva capa de delivery
   - Agregar CLI: agregar nuevos handlers CLI

🚫 ERRORES COMUNES QUE ESTA ARQUITECTURA EVITA:

❌ Lógica de negocio en controllers/handlers
❌ Dependencias directas a base de datos en controllers
❌ Código fuertemente acoplado
❌ Dificultad para testear
❌ Cambios que rompen múltiples partes del sistema

🎓 PREGUNTA PARA REFLEXIONAR:
¿Qué pasaría si quisiéramos:
- Cambiar de REST API a GraphQL?
- Agregar una interfaz CLI además de HTTP?
- Cambiar de memoria a PostgreSQL?
- Agregar cache Redis?
- Implementar autenticación JWT?

Respuesta: ¡Podríamos hacer todos estos cambios sin tocar
la lógica de negocio (casos de uso) ni las entidades!

🌟 ¡ESO ES EL PODER DE CLEAN ARCHITECTURE!
*/

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
