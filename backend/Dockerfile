FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar módulos
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o finanzas-api main.go

# Imagen de producción (multistage)
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binario compilado
COPY --from=builder /app/finanzas-api .
COPY DATABASE_SCHEMA.sql .

# Exponer puerto
EXPOSE 8080

# Comando de inicio
CMD ["./finanzas-api"]
