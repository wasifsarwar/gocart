# 🛒 GoCart - Microservices E-commerce API

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

---

## 🏗️ **Architecture Overview**

GoCart implements a **microservices architecture** with the following services:

### ✅ **Implemented Services**

| Service | Port | Status | Description |
|---------|------|--------|-------------|
| **Product Service** | `:8080` | ✅ Complete | Product catalog, inventory, CRUD operations |
| **User Service** | `:8081` | ✅ Complete | User management, authentication, profiles |

### 🚧 **Planned Services**
- **Order Service** - Order processing and management
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

### **Run with Coverage**
```bash
go test -v ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### **Integration Tests**
```bash
# Product Service Integration Tests
cd internal/product-service && go test -v ./...

# User Service Integration Tests  
cd internal/user-service && go test -v ./...
```

### **Test Database Isolation**
Each integration test creates its own isolated PostgreSQL database:
- ✅ **No test interference** - Each test runs in isolation
- ✅ **Parallel execution** - Tests can run concurrently
- ✅ **Automatic cleanup** - Databases are dropped after tests

---

## 🔧 **Technology Stack**

### **Backend**
- **Language**: Go 1.23
- **Framework**: Gorilla Mux (HTTP routing)
- **Database**: PostgreSQL + GORM ORM
- **Testing**: Go testing + testify
- **Documentation**: OpenAPI 3.0.3 + Swagger UI

### **Infrastructure**
- **Containerization**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Database Migrations**: GORM AutoMigrate
- **Code Coverage**: Go built-in tools + Codecov

### **Development**
- **Project Layout**: Standard Go project layout
- **Dependency Management**: Go Modules
- **Code Quality**: gofmt, go vet, golint
- **Version Control**: Git with conventional commits

---

## 🌟 **Key Features**

### **🏗️ Architecture**
- ✅ **Microservices** - Independent, scalable services
- ✅ **RESTful APIs** - Standard HTTP/JSON interfaces
- ✅ **Concurrent Execution** - Services run simultaneously

### **🔒 Data Management**
- ✅ **PostgreSQL** - Robust relational database
- ✅ **GORM Integration** - Type-safe database operations
- ✅ **UUID Primary Keys** - Globally unique identifiers
- ✅ **Automatic Migrations** - Schema management

### **📋 Testing Strategy**
- ✅ **Unit Tests** - Business logic validation
- ✅ **Integration Tests** - End-to-end workflows
- ✅ **CI/CD Pipeline** - Automated testing

---

## 🔄 **CI/CD Pipeline**

1. **🔍 Code Quality** - Runs linting and formatting checks
2. **🏗️ Build** - Compiles all services
3. **🧪 Test** - Executes unit and integration tests
4. **📊 Coverage** - Generates coverage reports
5. **📋 Artifacts** - Uploads test results and coverage

---

## 🌐 **Environment Configuration**

### **Database Configuration**
```bash
# Development
export DB_HOST=localhost
export DB_USER=admin
export DB_PASSWORD=admin
export DB_NAME=gocart_db
export DB_PORT=5432

# Testing (automatically handled)
export TEST_DB_HOST=localhost
export TEST_DB_USER=admin
export TEST_DB_PASSWORD=admin
export TEST_DB_NAME=gocart_db
export TEST_DB_PORT=5432
```

### **Service Ports**
```bash
export PRODUCT_SERVICE_PORT=8080
export USER_SERVICE_PORT=8081
```

---

## 🚧 **Roadmap**

### **Phase 1: Foundation** ✅ *Complete*
- [x] Microservices architecture
- [x] Product & User services
- [x] PostgreSQL integration
- [x] Comprehensive testing
- [x] CI/CD pipeline
- [x] API documentation

### **Phase 2: Core Features** 🚧 *In Progress*
- [ ] Order Service implementation
- [ ] Payment processing simulation
- [ ] JWT authentication
- [ ] Role-based access control

### **Phase 3: Advanced Features** 📋 *Planned*
- [ ] Redis caching layer
- [ ] Rate limiting
- [ ] Message queue integration
- [ ] Circuit breaker pattern
- [ ] Monitoring & metrics

### **Phase 4: Production Ready** 🎯 *Future*
- [ ] Load balancing
- [ ] Service mesh (Istio)
- [ ] Distributed tracing
- [ ] Kubernetes deployment

---

## 📞 **Contact & Support**

- **GitHub**: [@wasifsarwar](https://github.com/wasifsarwar)
- **Issues**: [GitHub Issues](https://github.com/wasifsarwar/gocart/issues)
- **Documentation**: [API Docs](https://wasifsarwar.github.io/gocart/)

---

<div align="center">

[API Documentation](https://wasifsarwar.github.io/gocart/) • [Report Bug](https://github.com/wasifsarwar/gocart/issues) • [Request Feature](https://github.com/wasifsarwar/gocart/issues)

</div>

