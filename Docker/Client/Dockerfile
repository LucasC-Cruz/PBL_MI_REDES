FROM golang:1.21-alpine
WORKDIR /app
COPY . .
WORKDIR /app/Client
RUN go build -o cliente .
CMD ["./cliente"]