// Package postgresql contiene las implementaciones de repositorios usando PostgreSQL
// Este es un EJEMPLO de cómo implementar un repositorio real con base de datos
//
// 📚 ¿Cuándo usar esta implementación?
// - Cuando necesites persistencia real (los datos sobreviven al reinicio)
// - En entornos de producción
// - Cuando necesites features avanzadas de DB (transacciones, índices, etc.)
//
// 🔧 Para usar esta implementación:
// 1. Instalar driver PostgreSQL: go get github.com/lib/pq
// 2. Cambiar en main.go: memory.NewInMemoryBookRepository() → postgresql.NewPostgresBookRepository(db)
// 3. Ejecutar docker-compose up para levantar PostgreSQL
//
// 💡 NOTA: Este archivo es solo un EJEMPLO educativo
// En una implementación real, usarías un ORM como GORM o un query builder
package postgresql

import (
	"database/sql"
	"errors"
	"go-book-clean-architecture-api/internal/domain"
	"go-book-clean-architecture-api/internal/repository"
	// _ "github.com/lib/pq" // Driver PostgreSQL - Comentado porque no está instalado
)

// PostgresBookRepository implementa BookRepository usando PostgreSQL
//
// 🗃️ Diferencias con la implementación en memoria:
// ✅ Los datos persisten entre reinicios
// ✅ Soporte para transacciones
// ✅ Mejor performance con grandes volúmenes de datos
// ✅ Índices para búsquedas rápidas
// ❌ Más complejo de configurar
// ❌ Requiere base de datos externa
type PostgresBookRepository struct {
	db *sql.DB // Conexión a PostgreSQL
}

// NewPostgresBookRepository crea una nueva instancia del repositorio PostgreSQL
//
// 🔧 Ejemplo de uso:
// db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
// if err != nil { log.Fatal(err) }
// bookRepo := postgresql.NewPostgresBookRepository(db)
func NewPostgresBookRepository(db *sql.DB) repository.BookRepository {
	return &PostgresBookRepository{
		db: db,
	}
}

// Create almacena un nuevo libro en PostgreSQL
func (r *PostgresBookRepository) Create(book *domain.Book) (*domain.Book, error) {
	query := `
		INSERT INTO books (id, title, author) 
		VALUES ($1, $2, $3) 
		RETURNING id, title, author, created_at`

	var createdBook domain.Book
	var createdAt string // Para capturar created_at si necesitas

	err := r.db.QueryRow(query, book.ID, book.Title, book.Author).Scan(
		&createdBook.ID,
		&createdBook.Title,
		&createdBook.Author,
		&createdAt,
	)

	if err != nil {
		return nil, err
	}

	return &createdBook, nil
}

// GetByID busca un libro por su ID en PostgreSQL
func (r *PostgresBookRepository) GetByID(id string) (*domain.Book, error) {
	query := `SELECT id, title, author FROM books WHERE id = $1`

	var book domain.Book
	err := r.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("libro no encontrado")
		}
		return nil, err
	}

	return &book, nil
}

// GetAll retorna todos los libros desde PostgreSQL
func (r *PostgresBookRepository) GetAll() ([]*domain.Book, error) {
	query := `SELECT id, title, author FROM books ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		var book domain.Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}

// Update modifica un libro existente en PostgreSQL
func (r *PostgresBookRepository) Update(book *domain.Book) (*domain.Book, error) {
	query := `
		UPDATE books 
		SET title = $2, author = $3, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1 
		RETURNING id, title, author`

	var updatedBook domain.Book
	err := r.db.QueryRow(query, book.ID, book.Title, book.Author).Scan(
		&updatedBook.ID,
		&updatedBook.Title,
		&updatedBook.Author,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("libro no encontrado")
		}
		return nil, err
	}

	return &updatedBook, nil
}

// Delete elimina un libro por su ID en PostgreSQL
func (r *PostgresBookRepository) Delete(id string) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("libro no encontrado")
	}

	return nil
}

// PostgresUserRepository implementa UserRepository usando PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository crea una nueva instancia del repositorio PostgreSQL para usuarios
func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// Create almacena un nuevo usuario en PostgreSQL
func (r *PostgresUserRepository) Create(user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (id, name, email) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, email`

	var createdUser domain.User
	err := r.db.QueryRow(query, user.ID, user.Name, user.Email).Scan(
		&createdUser.ID,
		&createdUser.Name,
		&createdUser.Email,
	)

	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

// GetByID busca un usuario por su ID en PostgreSQL
func (r *PostgresUserRepository) GetByID(id string) (*domain.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`

	var user domain.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &user, nil
}

// GetAll retorna todos los usuarios desde PostgreSQL
func (r *PostgresUserRepository) GetAll() ([]*domain.User, error) {
	query := `SELECT id, name, email FROM users ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// Update modifica un usuario existente en PostgreSQL
func (r *PostgresUserRepository) Update(user *domain.User) (*domain.User, error) {
	query := `
		UPDATE users 
		SET name = $2, email = $3, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1 
		RETURNING id, name, email`

	var updatedUser domain.User
	err := r.db.QueryRow(query, user.ID, user.Name, user.Email).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &updatedUser, nil
}

// Delete elimina un usuario por su ID en PostgreSQL
func (r *PostgresUserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("usuario no encontrado")
	}

	return nil
}

// 🔧 PARA USAR ESTA IMPLEMENTACIÓN EN MAIN.GO:
//
// import (
//     "database/sql"
//     "go-book-clean-architecture-api/internal/infrastructure/postgresql"
//     _ "github.com/lib/pq"
// )
//
// func main() {
//     // Conectar a PostgreSQL
//     db, err := sql.Open("postgres",
//         "postgres://postgres:password@localhost/cleanarch?sslmode=disable")
//     if err != nil {
//         log.Fatal("Error conectando a PostgreSQL:", err)
//     }
//     defer db.Close()
//
//     // Crear repositorios PostgreSQL
//     bookRepo := postgresql.NewPostgresBookRepository(db)
//     userRepo := postgresql.NewPostgresUserRepository(db)
//
//     // El resto del código permanece igual...
// }

// 💡 VENTAJAS de esta implementación:
// ✅ Persistencia real de datos
// ✅ Transacciones ACID
// ✅ Consultas complejas con SQL
// ✅ Índices para mejor performance
// ✅ Backup y recovery
// ✅ Concurrencia avanzada
//
// 🔧 MEJORAS QUE PODRÍAS AGREGAR:
// - Connection pooling
// - Manejo de transacciones
// - Logging de queries
// - Health checks de la DB
// - Métricas de performance
// - Retry logic
//
// 🎓 PREGUNTA PARA REFLEXIONAR:
// ¿Tuvimos que cambiar nuestros casos de uso para usar PostgreSQL?
// Respuesta: ¡NO! Ese es el poder de Clean Architecture.
