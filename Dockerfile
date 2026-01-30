# Etapa 1: Compilación
FROM golang:1.23-alpine AS builder

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos go.mod y go.sum antes de copiar el código fuente
COPY go.mod ./

# Descarga las dependencias (usando go mod)
RUN go mod download

# Copia los archivos del proyecto al contenedor
COPY . .

# Descarga las dependencias (si es necesario)
RUN go mod tidy

# Compila el binario de la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o artdo-web

# Etapa 2: Imagen ligera para producción
FROM alpine:latest

# Instalar certificados CA para conexiones HTTPS (Gmail, Google API)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia el binario compilado desde la etapa de construcción
COPY --from=builder /app/artdo-web .

# Copia los archivos de configuración y recursos
COPY --from=builder /app/config.json .
COPY --from=builder /app/artdotech-core.css .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/locales ./locales

# Expone el puerto en el que la aplicación escucha
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./artdo-web"]
