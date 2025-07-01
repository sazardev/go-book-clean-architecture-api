// Package domain contiene las entidades principales de nuestro negocio
// Esta es la CAPA MÁS IMPORTANTE de Clean Architecture
//
// 🎯 REGLAS DE ORO del Dominio:
// ✅ NO debe depender de NADA externo (frameworks, DBs, HTTP, etc.)
// ✅ Contiene la lógica de negocio más crítica e importante
// ✅ Las otras capas pueden depender de esta, pero esta NO depende de otras
// ✅ Si cambias de framework web, base de datos, etc., esta capa NO cambia
//
// 🤔 ¿Por qué separamos las entidades?
// - Son el corazón de nuestro negocio
// - Representan conceptos del mundo real
// - Definen qué datos son importantes para nuestro sistema
package domain

// Book representa la entidad principal de nuestro dominio de libros
//
// 📖 ¿Qué es una entidad en Clean Architecture?
// - Es un objeto que tiene identidad única (ID)
// - Contiene datos y comportamientos relacionados con un concepto del negocio
// - En este caso simple, solo contiene datos, pero podría tener métodos de validación
//
// 🎯 Ejemplo de método que podríamos agregar:
// func (b *Book) IsValid() bool {
//     return b.Title != "" && b.Author != ""
// }
type Book struct {
	ID     string `json:"id"`     // Identificador único del libro
	Title  string `json:"title"`  // Título del libro
	Author string `json:"author"` // Autor del libro
}

// User representa la entidad de usuario en nuestro dominio
//
// 👤 ¿Por qué tenemos User además de Book?
// - Para demostrar cómo manejar múltiples entidades en Clean Architecture
// - En aplicaciones reales, tendrías decenas de entidades (Customer, Order, Product, etc.)
// - Cada entidad se maneja con el mismo patrón
//
// 🔍 Nota: Mantenemos las entidades simples y enfocadas en una sola responsabilidad
type User struct {
	ID    string `json:"id"`    // Identificador único del usuario
	Name  string `json:"name"`  // Nombre del usuario
	Email string `json:"email"` // Email del usuario
}

// 💡 CONSEJOS PARA PRINCIPIANTES:
//
// 1. 📝 Las entidades deben ser simples y reflejar conceptos del mundo real
// 2. 🚫 NO incluyas aquí lógica de HTTP, base de datos, o frameworks
// 3. ✅ SÍ puedes agregar métodos de validación o comportamientos del negocio
// 4. 🎯 Si tienes dudas, pregúntate: "¿Esto es parte del negocio o de la tecnología?"
//
// 🌟 EJEMPLO DE LO QUE PODRÍAS AGREGAR:
// - Métodos de validación: IsValidEmail(), HasRequiredFields()
// - Comportamientos del negocio: CalculateAge(), FormatFullName()
// - Constantes del dominio: MaxTitleLength, ValidEmailRegex
//
// 🚫 EJEMPLO DE LO QUE NO DEBES AGREGAR:
// - Anotaciones de base de datos: @Table, @Column
// - Lógica HTTP: ParseFromJSON(), ToHTTPResponse()
// - Dependencias externas: logging, frameworks, etc.
