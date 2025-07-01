# Dockerfile para la aplicación Go Clean Architecture

# Etapa 1: Build
FROM golang:1.21-alpine AS builder

# Instalar herramientas necesarias
RUN apk add --no-cache git

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Etapa 2: Runtime
FROM alpine:latest

# Instalar certificados SSL
RUN apk --no-cache add ca-certificates

# Crear directorio de trabajo
WORKDIR /root/

# Copiar el binario desde la etapa de build
COPY --from=builder /app/main .

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
