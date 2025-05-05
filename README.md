# Go REST API

A RESTful API built with Go, Gin framework, and MongoDB.

## Features

- User authentication (signup, login) with JWT
- Product management (CRUD operations)
- MongoDB database integration
- Middleware for authentication and request handling
- Environment-based configuration

## Tech Stack

- [Go](https://golang.org/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [MongoDB](https://www.mongodb.com/)
- [JWT Authentication](https://github.com/golang-jwt/jwt)
- [Air](https://github.com/cosmtrek/air) (live reload)

## Project Structure

```markdown
├── config/         # Application configuration
├── controllers/    # API controllers (auth, product)
├── database/       # Database connection and operations
├── middleware/     # Custom middleware
├── models/         # Data models
├── routers/        # API routes
├── utils/          # Utility functions
├── .env            # Environment variables
├── air.toml        # Air configuration for live reload
├── go.mod          # Go module dependencies
├── go.sum          # Go module checksum
├── main.go         # Application entry point
└── Makefile        # Build and development commands
```

## Prerequisites

- Go 1.16 or higher
- MongoDB database

## Getting Started

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/go_rest_api.git
   cd go_rest_api
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Configure environment variables:
   Create a `.env` file in the root directory with the following variables:

   ```bash
   PORT=8000
   MONGODB_URI=your_mongodb_connection_string
   DATABASE_NAME=go_rest_api
   GIN_MODE=debug  # 'debug' or 'release'
   JWT_SECRET_KEY=your_secret_key
   ```

### Running the Application

#### Development mode (with live reload)

```bash
make dev
```

#### Standard run

```bash
make run
```

#### Build and run

```bash
make build
make start
```

## API Endpoints

### Authentication

- `POST /api/auth/signup` - Register a new user
- `POST /api/auth/login` - Login and get JWT token

### Products

- `GET /api/products` - Get all products
- `GET /api/products/:id` - Get a specific product
- `POST /api/products` - Create a new product
- `PUT /api/products/:id` - Update a product
- `DELETE /api/products/:id` - Delete a product

## License

This project is licensed under the MIT License.
