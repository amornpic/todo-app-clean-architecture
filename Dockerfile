FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o todo-app main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/todo-app .
# Copy the .env file
COPY .env .  
EXPOSE 8080
CMD ["./todo-app"]