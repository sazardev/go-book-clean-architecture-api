# 📁 Estructura Completa del Proyecto

```
go-book-clean-architecture-api/
│
├── 📄 README.md                           # Guía principal y documentación
├── 📄 GUIDE.md                            # Guía paso a paso detallada
├── 📄 DEPLOYMENT.md                       # Guía de despliegue y características avanzadas
├── 📄 api_examples.http                   # Ejemplos de peticiones HTTP para probar
├── 📄 go.mod                              # Definición del módulo Go y dependencias
├── 📄 go.sum                              # Checksums de dependencias
├── 📄 .gitignore                          # Archivos ignorados por Git
├── 📄 Dockerfile                          # Imagen Docker para la aplicación
├── 📄 docker-compose.yml                  # Servicios Docker (app + postgres + redis)
│
├── 📁 cmd/                                # Puntos de entrada de la aplicación
│   └── 📁 server/
│       └── 📄 main.go                     # 🚀 Archivo principal - Dependency Injection
│
├── 📁 internal/                           # Código interno de la aplicación
│   │
│   ├── 📁 domain/                         # 🏛️ CAPA DE DOMINIO
│   │   └── 📄 book.go                     # Entidades: Book y User
│   │
│   ├── 📁 repository/                     # 📋 CONTRATOS/INTERFACES
│   │   └── 📄 book_repository.go          # BookRepository y UserRepository interfaces
│   │
│   ├── 📁 usecase/                        # 🧠 CAPA DE APLICACIÓN/CASOS DE USO
│   │   ├── 📄 book_usecause.go            # BookUseCase y UserUseCase
│   │   └── 📁 test/
│   │       └── 📄 book_usecase_test.go    # 🧪 Tests de casos de uso
│   │
│   ├── 📁 infrastructure/                 # 💾 CAPA DE INFRAESTRUCTURA
│   │   ├── 📁 memory/
│   │   │   └── 📄 book_repository.go      # Implementación en memoria
│   │   └── 📁 postgresql/
│   │       └── 📄 book_repository.go      # 🗃️ Implementación PostgreSQL (ejemplo)
│   │
│   ├── 📁 delivery/                       # 🌐 CAPA DE DELIVERY/INTERFAZ
│   │   └── 📁 http/
│   │       └── 📄 book_handler.go         # BookHandler y UserHandler HTTP
│   │
│   └── 📁 routes/                         # 🛣️ CONFIGURACIÓN DE RUTAS
│       └── 📄 book_routes.go              # Definición de todas las rutas
│
├── 📁 scripts/                            # Scripts de base de datos
│   └── 📄 init.sql                        # Inicialización PostgreSQL
│
└── 📁 .vscode/                            # Configuración VS Code
    └── 📄 tasks.json                      # Tareas predefinidas (ejecutar, compilar, testear)
```

## 🔄 Flujo de Datos entre Capas

```
📱 Cliente HTTP
    ↓
🌐 HTTP Handler (Delivery Layer)
    ↓
🧠 Use Case (Application Layer)
    ↓
📋 Repository Interface (Application Layer)
    ↓
💾 Repository Implementation (Infrastructure Layer)
    ↓
🗄️ Base de Datos / Memoria
```

## 📚 Archivos por Capa de Clean Architecture

### 🏛️ Domain Layer (Independiente de todo)
- `internal/domain/book.go` - Entidades Book y User

### 🧠 Application Layer (Casos de uso + Interfaces)
- `internal/usecase/book_usecause.go` - Lógica de negocio
- `internal/repository/book_repository.go` - Contratos/Interfaces
- `internal/usecase/test/book_usecase_test.go` - Tests de casos de uso

### 💾 Infrastructure Layer (Implementaciones concretas)
- `internal/infrastructure/memory/book_repository.go` - Repositorio en memoria
- `internal/infrastructure/postgresql/book_repository.go` - Repositorio PostgreSQL

### 🌐 Interface/Delivery Layer (Adaptadores externos)
- `internal/delivery/http/book_handler.go` - Handlers HTTP
- `internal/routes/book_routes.go` - Configuración de rutas

### ⚙️ Main/Composition Root
- `cmd/server/main.go` - Inyección de dependencias y configuración

## 🎯 Endpoints Disponibles

### 📖 Gestión de Libros
- `POST /api/books` - Crear libro
- `GET /api/books` - Obtener todos los libros
- `GET /api/books/:id` - Obtener libro por ID
- `PUT /api/books/:id` - Actualizar libro
- `DELETE /api/books/:id` - Eliminar libro

### 👤 Gestión de Usuarios
- `POST /api/users` - Crear usuario
- `GET /api/users` - Obtener todos los usuarios
- `GET /api/users/:id` - Obtener usuario por ID
- `PUT /api/users/:id` - Actualizar usuario
- `DELETE /api/users/:id` - Eliminar usuario

### 🔍 Otros
- `GET /health` - Health check

## 🚀 Cómo empezar

1. **Leer documentación:**
   - Empezar con `README.md`
   - Continuar con `GUIDE.md` para entender paso a paso
   - Ver `DEPLOYMENT.md` para características avanzadas

2. **Ejecutar la aplicación:**
   ```bash
   go run cmd/server/main.go
   ```

3. **Probar la API:**
   - Abrir `api_examples.http` en VS Code
   - Instalar extensión "REST Client"
   - Hacer clic en "Send Request" en los ejemplos

4. **Ejecutar tests:**
   ```bash
   go test ./internal/usecase/test -v
   ```

5. **Experimentar:**
   - Modificar entidades en `internal/domain/`
   - Agregar nuevos casos de uso en `internal/usecase/`
   - Crear nuevos endpoints en `internal/delivery/http/`

## 💡 Conceptos Clave Implementados

✅ **Clean Architecture** - Separación estricta en 4 capas
✅ **Dependency Injection** - Inyección manual en `main.go`
✅ **Repository Pattern** - Abstracción de persistencia
✅ **Use Cases** - Lógica de negocio encapsulada
✅ **Testing** - Tests unitarios de casos de uso
✅ **Docker** - Contenedorización completa
✅ **PostgreSQL** - Ejemplo de repositorio real
✅ **Documentación** - Comentarios educativos extensos
✅ **VS Code Integration** - Tareas predefinidas
✅ **HTTP Examples** - Ejemplos listos para usar

## 🎓 Para Aprender Más

- Estudiar el flujo de datos entre capas
- Entender cómo se inyectan las dependencias
- Practicar agregando nuevas entidades (ej: Author, Category)
- Experimentar cambiando de memoria a PostgreSQL
- Crear tests para handlers y repositorios
- Implementar características avanzadas del `DEPLOYMENT.md`

---

**¡Esta es tu plantilla completa de Clean Architecture en Go! 🚀**
