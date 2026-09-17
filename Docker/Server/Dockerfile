FROM golang:1.21-alpine
WORKDIR /app
# Copia todo o projeto (incluindo o go.mod da raiz e a pasta Server)
COPY . .
WORKDIR /app/Server
RUN go build -o servidor servidor.go
EXPOSE 6742
CMD ["./servidor", "6742"]