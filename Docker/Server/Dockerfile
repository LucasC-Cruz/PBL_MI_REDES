FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copia todo o projeto para resolver dependências internas (Lgcc, Model, etc.)
COPY . .

# Compila o binário do servidor
RUN go build -o server ./Server/servidor.go

FROM alpine:latest
WORKDIR /app

# Copia o binário gerado
COPY --from=builder /app/server .

# Mantém a estrutura 'PBL_MI_REDES/data' esperada pelas chamadas os.ReadFile no servidor
COPY --from=builder /app/data ./PBL_MI_REDES/data

EXPOSE 6742

# Executa o servidor passando a porta 6742 via parâmetro
CMD ["./server", "6742"]