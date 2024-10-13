# Messaging App

This is a Go-based service designed to handle scheduled message sending, using GORM for database management, Redis for caching, and HTTP requests to send messages. The application allows you to send messages periodically to a webhook URL, storing the messages in a SQLite database and sent messsages in Redis.

## Features

- **Message Scheduling**: Messages are scheduled and sent at regular intervals.
- **Database Integration**: Uses SQLite (via GORM) to store messages and track their status.
- **Redis Caching**: Utilizes Redis for caching message IDs and timestamps after being sent.
- **State Management**: Implements state design pattern to manage service states (started/stopped).

## Prerequisites

- Go 1.16+
- Docker (optional but recommended)
- Redis server (local or Dockerized)

## Installation

1. **Install Dependencies:**

   If you're using Go modules, dependencies should automatically be handled when building or running the application.

   ```bash
   go mod tidy
   ```

2. **Set Up Environment Variables (Required!!):**

   You can create a `.env` file or set environment variables directly to customize the configuration:

   ```env
   DATABASE_URL=messages.db
   REDIS_ADDR=localhost:6379
   WEBHOOK_URL=https://webhook.site/your-endpoint
   MESSAGING_PERIOD=120 # in seconds
   PORT=8080
   ```

3. **Run the Application:**

   ```bash
   go run cmd/myapp/main.go
   ```

## Configuration

The configuration variables are stored in the `config/config.go` file and can be loaded from environment variables or `.env` file. Here are the default configuration options:

- `WEBHOOK_URL`: URL to send messages to (default: `https://webhook.site`)
- `PORT`: Port on which the application runs (default: `8080`)
- `DATABASE_URL`: Path to the SQLite database file (default: `messages.db`)
- `REDIS_ADDR`: Redis server address (default: `localhost:6379`)
- `MESSAGING_PERIOD`: Period between message sends (default: `120` seconds)


## Running with Docker

You can use Docker and Docker Compose to easily run the application with Redis.

1. **Run the Application with Docker Compose:**

   ```bash
   docker-compose up --build
   ```

This will start both the Go application and Redis in separate containers.

## API Endpoints

The application exposes the following API endpoints:

- `GET /startstop?action=` Start or stop the messaging service by passing `start` or `stop` as the `action` query parameter.

- `GET /get-sent-messages` Getting all messages that are sent.

Also you can use the swagger documentation to see the API documentation by visiting `/swagger` endpoint.

## Testing
The application includes unit tests, which can be run using:

```bash
go test ./internal/service
```

## Folder Structure

- `cmd/`: Contains the main entry point for the application.
- `service/`: Business logic, including message scheduling and state management.
- `db/`: Database models and operations (using GORM).
- `cache/`: Redis client and caching logic.
- `config/`: Configuration handling.

