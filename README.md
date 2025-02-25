# Todo App

A simple Todo application built with Go, following clean architecture principles.

## Prerequisites

- Go 1.21 or higher
- Git
- Docker

## Project Structure
```
todo-app-clean-architecture/
├── domain/ # Enterprise business rules and entities
├── usecase/ # Application business rules
├── repository/ # Data layer implementations
├── api/ # Interface adapters (HTTP handlers)
```

## Getting Started

### 1. Clone the repository
`git clone https://github.com/amornpic/todo-app-clean-architecture.git`

`cd todo-app-clean-architecture`

### 2. Set up environment variables
`cp .env.example .env`

### 3. Install dependencies
`go mod tidy`

### 4. Start project
Using Docker
`docker-compose up -d`

The server will start on `http://localhost:8080`

## API Documentation

API documentation is available via Swagger UI at:
http://localhost:8080/swagger/index.html

To update Swagger documentation after making API changes:
`swag init`

## Running Tests
`go test ./...`

## API Endpoints

- `POST /todos` - Create a new todo
- `GET /todos` - List all todos
- `PUT /todos/{id}` - Update a todo
- `DELETE /todos/{id}` - Delete a todo