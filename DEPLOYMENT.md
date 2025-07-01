# 🚀 Guía de Despliegue y Características Avanzadas

## 📋 Tabla de Contenidos

1. [Desarrollo Local](#desarrollo-local)
2. [Testing](#testing)
3. [Docker y Contenedores](#docker-y-contenedores)
4. [Base de Datos PostgreSQL](#base-de-datos-postgresql)
5. [Monitoreo y Logging](#monitoreo-y-logging)
6. [Deployment en Producción](#deployment-en-producción)
7. [Características Avanzadas](#características-avanzadas)

## 🛠️ Desarrollo Local

### Prerrequisitos
- Go 1.21 o superior
- Docker y Docker Compose (opcional)
- VS Code con extensión Go (recomendado)

### Configuración Inicial
```bash
# Clonar el repositorio
git clone <repo-url>
cd go-book-clean-architecture-api

# Instalar dependencias
go mod tidy

# Ejecutar la aplicación
go run cmd/server/main.go
```

### Usando VS Code Tasks
Este proyecto incluye tareas predefinidas en `.vscode/tasks.json`:

- **🚀 Ejecutar Servidor**: `Ctrl+Shift+P` → "Tasks: Run Task" → "🚀 Ejecutar Servidor"
- **🔨 Compilar**: Compila la aplicación en un ejecutable
- **🧪 Ejecutar Tests**: Ejecuta todos los tests del proyecto
- **📦 Instalar Dependencias**: Ejecuta `go mod tidy`

## 🧪 Testing

### Ejecutar Tests
```bash
# Todos los tests
go test ./...

# Tests con coverage
go test -cover ./...

# Tests específicos con verbose
go test -v ./internal/usecase/test

# Coverage detallado
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Estructura de Testing
```
internal/
├── usecase/
│   └── test/
│       └── book_usecase_test.go  # Tests de casos de uso
├── delivery/
│   └── http/
│       └── test/
│           └── book_handler_test.go  # Tests de handlers (a crear)
└── infrastructure/
    └── memory/
        └── test/
            └── book_repository_test.go  # Tests de repositorios (a crear)
```

### Creando Tests para Handlers
```go
// Ejemplo: internal/delivery/http/test/book_handler_test.go
func TestCreateBookHandler(t *testing.T) {
    // Arrange
    mockUseCase := &MockBookUseCase{}
    handler := http.NewBookHandler(mockUseCase)
    
    app := fiber.New()
    app.Post("/books", handler.CreateBook)
    
    body := `{"title": "Test Book", "author": "Test Author"}`
    req := httptest.NewRequest("POST", "/books", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    // Act
    resp, err := app.Test(req)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
}
```

## 🐳 Docker y Contenedores

### Desarrollo con Docker

#### Solo la aplicación
```bash
# Construir imagen
docker build -t clean-arch-api .

# Ejecutar contenedor
docker run -p 8080:8080 clean-arch-api
```

#### Ambiente completo con Docker Compose
```bash
# Levantar todos los servicios (app + postgres + redis)
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down

# Reconstruir imágenes
docker-compose up --build
```

### Servicios incluidos en docker-compose.yml
- **app**: La aplicación Go Clean Architecture
- **postgres**: Base de datos PostgreSQL con datos iniciales
- **redis**: Cache Redis (para futuras implementaciones)

### Configuración de Producción
```dockerfile
# Dockerfile optimizado para producción
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

## 🗄️ Base de Datos PostgreSQL

### Cambiar de Memoria a PostgreSQL

1. **Instalar driver PostgreSQL**
```bash
go get github.com/lib/pq
```

2. **Modificar main.go**
```go
import (
    "database/sql"
    "go-book-clean-architecture-api/internal/infrastructure/postgresql"
    _ "github.com/lib/pq"
)

func main() {
    // Conectar a PostgreSQL
    db, err := sql.Open("postgres", 
        "postgres://postgres:password@localhost/cleanarch?sslmode=disable")
    if err != nil {
        log.Fatal("Error conectando a PostgreSQL:", err)
    }
    defer db.Close()

    // Usar repositorios PostgreSQL
    bookRepo := postgresql.NewPostgresBookRepository(db)
    userRepo := postgresql.NewPostgresUserRepository(db)
    
    // El resto permanece igual...
}
```

3. **Inicializar base de datos**
```bash
# Con Docker Compose (automático)
docker-compose up postgres

# Manual
psql -h localhost -U postgres -d cleanarch -f scripts/init.sql
```

### Migraciones de Base de Datos

Para entornos de producción, considera usar herramientas como:
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)
- [Atlas](https://atlasgo.io/)

Ejemplo con golang-migrate:
```bash
# Instalar migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Crear migración
migrate create -ext sql -dir migrations -seq create_books_table

# Ejecutar migraciones
migrate -path migrations -database "postgres://postgres:password@localhost/cleanarch?sslmode=disable" up
```

## 📊 Monitoreo y Logging

### Logging Estructurado
```go
// Ejemplo con logrus
import "github.com/sirupsen/logrus"

func (uc *BookUseCase) CreateBook(title, author string) (*domain.Book, error) {
    logrus.WithFields(logrus.Fields{
        "title":  title,
        "author": author,
    }).Info("Creating new book")
    
    // ... resto de la lógica
    
    logrus.WithField("bookID", book.ID).Info("Book created successfully")
    return book, nil
}
```

### Métricas con Prometheus
```go
// Ejemplo de métricas
import "github.com/prometheus/client_golang/prometheus"

var (
    booksCreated = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "books_created_total",
            Help: "Total number of books created",
        },
        []string{"author"},
    )
)

func (uc *BookUseCase) CreateBook(title, author string) (*domain.Book, error) {
    book, err := uc.bookRepo.Create(&domain.Book{...})
    if err == nil {
        booksCreated.WithLabelValues(author).Inc()
    }
    return book, err
}
```

### Health Checks Avanzados
```go
// internal/delivery/http/health_handler.go
func (h *HealthHandler) DetailedHealthCheck(c *fiber.Ctx) error {
    health := map[string]interface{}{
        "status": "OK",
        "timestamp": time.Now(),
        "version": "1.0.0",
        "dependencies": map[string]string{
            "database": h.checkDatabase(),
            "redis":    h.checkRedis(),
        },
    }
    return c.JSON(health)
}
```

## 🚀 Deployment en Producción

### Kubernetes
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: clean-arch-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: clean-arch-api
  template:
    metadata:
      labels:
        app: clean-arch-api
    spec:
      containers:
      - name: api
        image: your-registry/clean-arch-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "postgres-service"
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
```

### Variables de Entorno
```go
// internal/config/config.go
type Config struct {
    Port       string `env:"PORT" envDefault:"8080"`
    DBHost     string `env:"DB_HOST" envDefault:"localhost"`
    DBUser     string `env:"DB_USER" envDefault:"postgres"`
    DBPassword string `env:"DB_PASSWORD" envDefault:"password"`
    DBName     string `env:"DB_NAME" envDefault:"cleanarch"`
    RedisURL   string `env:"REDIS_URL" envDefault:"redis://localhost:6379"`
}
```

### CI/CD con GitHub Actions
```yaml
# .github/workflows/ci.yml
name: CI/CD
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: 1.21
    - run: go mod tidy
    - run: go test ./...
    - run: go build ./cmd/server

  docker:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
    - uses: actions/checkout@v3
    - uses: docker/build-push-action@v3
      with:
        push: true
        tags: your-registry/clean-arch-api:latest
```

## 🌟 Características Avanzadas

### Autenticación JWT
```go
// internal/middleware/auth.go
func JWTMiddleware() fiber.Handler {
    return jwtware.New(jwtware.Config{
        SigningKey: []byte("your-secret-key"),
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            return c.Status(401).JSON(fiber.Map{
                "error": "Unauthorized",
            })
        },
    })
}

// En las rutas
books.Use(middleware.JWTMiddleware())
books.Post("/", bookHandler.CreateBook)
```

### Cache con Redis
```go
// internal/infrastructure/redis/cache.go
type RedisCache struct {
    client *redis.Client
}

func (c *RedisCache) Get(key string) (string, error) {
    return c.client.Get(context.Background(), key).Result()
}

func (c *RedisCache) Set(key, value string, ttl time.Duration) error {
    return c.client.Set(context.Background(), key, value, ttl).Err()
}

// En el caso de uso
func (uc *BookUseCase) GetBookByID(id string) (*domain.Book, error) {
    // Intentar obtener desde cache
    if cached, err := uc.cache.Get("book:" + id); err == nil {
        var book domain.Book
        json.Unmarshal([]byte(cached), &book)
        return &book, nil
    }
    
    // Si no está en cache, obtener desde DB
    book, err := uc.bookRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // Guardar en cache
    bookJSON, _ := json.Marshal(book)
    uc.cache.Set("book:"+id, string(bookJSON), 5*time.Minute)
    
    return book, nil
}
```

### Validación Avanzada
```go
// go get github.com/go-playground/validator/v10

type CreateBookRequest struct {
    Title  string `json:"title" validate:"required,min=1,max=255"`
    Author string `json:"author" validate:"required,min=1,max=255"`
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
    var req CreateBookRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid format"})
    }
    
    if err := validator.New().Struct(req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }
    
    // ... resto de la lógica
}
```

### Rate Limiting
```go
// internal/middleware/rate_limit.go
func RateLimitMiddleware() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        100,           // 100 requests
        Expiration: 1 * time.Minute, // per minute
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP() // Limit by IP
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(429).JSON(fiber.Map{
                "error": "Too many requests",
            })
        },
    })
}
```

### Documentación con Swagger
```go
// go get github.com/swaggo/swag/cmd/swag
// go get github.com/swaggo/fiber-swagger

// @title Clean Architecture API
// @version 1.0
// @description API RESTful siguiendo Clean Architecture
// @host localhost:8080
// @BasePath /api

// @Summary Crear un libro
// @Description Crea un nuevo libro en el sistema
// @Tags books
// @Accept json
// @Produce json
// @Param book body CreateBookRequest true "Datos del libro"
// @Success 201 {object} domain.Book
// @Failure 400 {object} map[string]string
// @Router /books [post]
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
    // ... implementación
}
```

### Event Sourcing (Avanzado)
```go
// internal/domain/events.go
type BookEvent interface {
    EventType() string
    AggregateID() string
    Timestamp() time.Time
}

type BookCreatedEvent struct {
    ID     string    `json:"id"`
    Title  string    `json:"title"`
    Author string    `json:"author"`
    Time   time.Time `json:"timestamp"`
}

func (e BookCreatedEvent) EventType() string { return "BookCreated" }
func (e BookCreatedEvent) AggregateID() string { return e.ID }
func (e BookCreatedEvent) Timestamp() time.Time { return e.Time }

// internal/infrastructure/eventstore/eventstore.go
type EventStore interface {
    SaveEvent(event domain.BookEvent) error
    GetEvents(aggregateID string) ([]domain.BookEvent, error)
}
```

## 📈 Métricas y Observabilidad

### Ejemplo completo con OpenTelemetry
```go
// internal/observability/tracing.go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func (uc *BookUseCase) CreateBook(ctx context.Context, title, author string) (*domain.Book, error) {
    tracer := otel.Tracer("book-usecase")
    ctx, span := tracer.Start(ctx, "CreateBook")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("book.title", title),
        attribute.String("book.author", author),
    )
    
    // ... lógica de negocio
    
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }
    
    span.SetAttributes(attribute.String("book.id", book.ID))
    return book, nil
}
```

## 🔒 Seguridad

### CORS, CSRF y Security Headers
```go
import (
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/csrf"
    "github.com/gofiber/fiber/v2/middleware/helmet"
)

app.Use(helmet.New()) // Security headers
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://yourdomain.com",
    AllowMethods: "GET,POST,PUT,DELETE",
    AllowHeaders: "Origin,Content-Type,Accept,Authorization",
}))
app.Use(csrf.New()) // CSRF protection
```

---

## 📚 Recursos Adicionales

- [Clean Architecture - Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Testing](https://go.dev/doc/tutorial/add-a-test)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Fiber Framework](https://docs.gofiber.io/)
- [PostgreSQL Go Driver](https://github.com/lib/pq)

¡Feliz coding y deployment! 🚀
