# 🛒 GoCart - Full-Stack E-commerce API

[![Go](https://github.com/wasifsarwar/gocart/workflows/Go/badge.svg)](https://github.com/wasifsarwar/gocart/actions)
[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![API Documentation](https://img.shields.io/badge/API-Documentation-blue)](https://wasifsarwar.github.io/gocart/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> A modern, scalable microservices-based e-commerce API built with Go, featuring comprehensive testing, CI/CD, and beautiful API documentation.

## 🔗 **Live API Documentation**

**[📖 Interactive API Explorer](https://wasifsarwar.github.io/gocart/)** - Test APIs directly in your browser

### Service-Specific Documentation
- **🛍️ Product API**: [Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/wasifsarwar/gocart/main/api/product/openapi.yaml) | [OpenAPI Spec](api/product/openapi.yaml)
- **👤 User API**: [Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/wasifsarwar/gocart/main/api/user/openapi.yaml) | [OpenAPI Spec](api/user/openapi.yaml)
- **🛒 Order API**: [Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/wasifsarwar/gocart/main/api/order-management/openapi.yaml) | [OpenAPI Spec](api/order-management/openapi.yaml)

---

## 🏗️ **Architecture Overview**

GoCart implements a **microservices architecture** with the following services:

### ✅ **Frontend Stack**
- **React 18** with TypeScript
- **React Router** for navigation  
- **Custom hooks** for state management
- **Responsive CSS** with mobile-first design
- **Real-time search & sorting**
- **Animated gopher branding** 🐹

### ✅ **Implemented Services**

| Service | Port | Status | Description |
|---------|------|--------|-------------|
| **Product Service** | `:8080` | ✅ Complete | Product catalog, inventory, CRUD operations |
| **User Service** | `:8081` | ✅ Complete | User management, authentication, profiles |
| **Order Service** | `:8082` | ✅ Complete | Order processing and management |

### 🚧 **Planned Services**
- **Payment Service** - Payment processing simulation
- **Notification Service** - Email/SMS notifications

---

## 🚀 **Quick Start**

### **Prerequisites**
- **Go 1.23+** 
- **PostgreSQL 13+**
- **Docker & Docker Compose** (optional)

### **1. Clone & Setup**
```bash
git clone https://github.com/wasifsarwar/gocart.git
cd gocart
go mod download
```

### **2. Database Setup**
```bash
# Start PostgreSQL with Docker
docker-compose up -d postgres

# Or use your local PostgreSQL
createdb gocart_db
```

```bash
# Shut down PostgreSQL instance
brew services stop postgresql
```

### **3. Run Services**

#### **Option A: Run All Services**
```bash
go run cmd/main.go
```
- Product Service: http://localhost:8080
- User Service: http://localhost:8081

#### **Option B: Docker Compose**
```bash
docker-compose up
```

```bash
# Full stack development
docker-compose up --build      # Start everything
docker-compose down           # Stop everything

# Frontend only
cd web && npm start           # Development server

# Backend only  
go run cmd/main.go           # Start Go services

# View logs
docker-compose logs frontend # Frontend logs
docker-compose logs app      # Backend logs
```
---

## 📊 **API Endpoints**

### **Product Service** (`localhost:8080`)
```http
GET    /products           # List all products
POST   /products           # Create product
GET    /products/{id}      # Get product by ID
PUT    /products/{id}      # Update product
DELETE /products/{id}      # Delete product
```

### **User Service** (`localhost:8081`)
```http
GET    /users              # List all users
POST   /users/register     # Register new user
GET    /users/{id}         # Get user by ID
PUT    /users/{id}         # Update user
DELETE /users/{id}         # Delete user
```

---

## 🏛️ **Project Structure**

```
gocart/
├── 📁 api/                          # OpenAPI specifications
│   ├── product/openapi.yaml         # Product service API spec
│   └── user/openapi.yaml            # User service API spec
├── 📁 cmd/
│   └── main.go                      # Main application entry point
├── 📁 docs/                         # Documentation and assets
│   └── index.html                   # Beautiful Swagger UI
├── 📁 internal/                     # Private application code
│   ├── 📁 product-service/
│   │   ├── handler/                 # HTTP handlers
│   │   ├── models/                  # Data models
│   │   ├── repository/              # Database layer
│   │   └── server/                  # Server setup
│   └── 📁 user-service/
│       ├── handler/                 # HTTP handlers
│       ├── models/                  # Data models
│       ├── repository/              # Database layer
│       └── server/                  # Server setup
├── 📁 pkg/                          # Shared utilities
│   ├── db/                          # Database connections
│   └── testutils/                   # Testing utilities
├── 📁 .github/workflows/            # CI/CD pipelines
├── docker-compose.yml               # Multi-service orchestration
├── Dockerfile                       # Container build instructions
├── go.mod                          # Go dependencies
└── README.md                       # You are here!
```

---

## 🧪 **Testing**

### **Run All Tests**
```bash
go test ./...
```
