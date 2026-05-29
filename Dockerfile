# ── Etapa 1: Build ───────────────────────────────────────────────────
FROM golang:1.22-alpine AS build

WORKDIR /app

# Descargar dependencias primero (cacheadas en Docker layers)
COPY go.mod go.sum ./
RUN go mod download

# Compilar la aplicación
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o task-queue-worker ./cmd/main.go

# ── Etapa 2: Runtime ─────────────────────────────────────────────────
FROM alpine:3.19

# Certificados SSL para SMTP
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copiar binario compilado
COPY --from=build /app/task-queue-worker .

# Copiar migraciones SQL
COPY --from=build /app/migrations ./migrations

EXPOSE ${PORT:-8080}

ENTRYPOINT ["./task-queue-worker"]
