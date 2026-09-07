# Go Backend Setup with PostgreSQL, GORM, Gin & Air

A modular, scalable, production-ready Go backend architecture built with **Gin Web Framework**, **GORM ORM**, **PostgreSQL**, and **Air** for live-reload development.

---

## 🚀 Features

- **Web Framework:** [Gin](https://github.com/gin-gonic/gin) - Fast & lightweight HTTP web framework.
- **ORM:** [GORM](https://gorm.io/) with PostgreSQL driver (`gorm.io/driver/postgres`).
- **Live Reload:** [Air](https://github.com/air-verse/air) (`.air.toml`) for hot reloading during development.
- **Swagger Documentation:** [Swagger UI](https://github.com/swaggo/gin-swagger) at `/swagger/index.html`.
- **Modular Architecture:** Clean separation of concerns with **Controller**, **Service**, and **Repository** layers.
- **Configuration:** Environment variable management via `.env` and `godotenv`.
- **Graceful Shutdown:** Safe server termination handling `SIGINT` / `SIGTERM`.
- **Standardized API Responses:** Clean JSON response helpers (`Success`, `Created`, `BadRequest`, `NotFound`, `InternalServerError`).

---

## 📂 Project Structure

```
├── config/
│   ├── config.go            # Environment variable configuration loader
│   └── database.go          # GORM PostgreSQL connection & pooling
├── internal/
│   ├── modules/
│   │   └── user/            # Modular User feature
│   │       ├── controller.go # Gin request handlers
│   │       ├── service.go    # Business logic layer
│   │       ├── repository.go # Database queries / persistence layer
│   │       ├── model.go      # GORM models & DTOs
│   │       └── routes.go     # Module route registration
│   └── routes/
│       └── routes.go        # Central router, CORS, and root "/" route
├── pkg/
│   └── response/
│       └── response.go      # Standardized JSON response helpers
├── .air.toml                # Air hot-reload configuration
├── .env.example             # Sample environment variables
├── .env                     # Local environment variables
├── .gitignore               # Git ignored files
├── go.mod                   # Go module definition
├── go.sum                   # Go dependencies checksum
├── main.go                  # Application entry point
└── readme.md                # Project documentation
```

---

## ⚙️ Getting Started

### 1. Prerequisites
- [Go (1.20+)](https://golang.org/dl/)
- [PostgreSQL](https://www.postgresql.org/)
- [Air](https://github.com/air-verse/air) (Optional, for hot reloading)

To install **Air** globally:
```bash
go install github.com/air-verse/air@latest
```

### 2. Environment Configuration
Copy `.env.example` to `.env` and update your PostgreSQL credentials:

```bash
cp .env.example .env
```

```env
PORT=8080
APP_ENV=development

# PostgreSQL Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=your_db_name
DB_SSLMODE=disable
DB_TIMEZONE=UTC
```

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Run Application

#### Using Air (Hot Reloading):
```bash
air
```

#### Standard Go Run:
```bash
go run main.go
```

---

## 📡 API Endpoints

### Root, Health & Swagger Docs
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/` | Root endpoint (`backend running successfully`) |
| `GET` | `/health` | Server health check |
| `GET` | `/swagger/index.html` | Interactive Swagger API Documentation UI |

### User Module (`/api/v1/users`)
| Method | Endpoint | Description | Payload Example |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/users` | Create new user | `{"name": "John Doe", "email": "john@example.com"}` |
| `GET` | `/api/v1/users` | List all users | N/A |
| `GET` | `/api/v1/users/:id` | Get user by ID | N/A |
| `PUT` | `/api/v1/users/:id` | Update user by ID | `{"name": "John Updated"}` |
| `DELETE` | `/api/v1/users/:id` | Soft delete user | N/A |

---

## 🏗️ Adding a New Module

To add a new module (e.g. `product`):
1. Create `internal/modules/product/` with:
   - `model.go` (GORM schema and request DTOs)
   - `repository.go` (GORM database operations)
   - `service.go` (Business logic)
   - `controller.go` (HTTP handlers)
   - `routes.go` (Route registration)
2. Register the module route group in [internal/routes/routes.go](file:///e:/go-backend-set-up-with-postgraySQL/internal/routes/routes.go).
3. Add the model to auto-migration in [main.go](file:///e:/go-backend-set-up-with-postgraySQL/main.go).
