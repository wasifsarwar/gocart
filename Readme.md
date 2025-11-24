# GoCart - Full-Stack E-commerce API

[![Go](https://github.com/wasifsarwar/gocart/workflows/Go/badge.svg)](https://github.com/wasifsarwar/gocart/actions)
[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![API Documentation](https://img.shields.io/badge/API-Documentation-blue)](https://wasifsarwar.github.io/gocart/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Netlify Status](https://api.netlify.com/api/v1/badges/578e84c6-0d06-4df9-b2fa-8c8993f653dd/deploy-status)](https://app.netlify.com/projects/gocartshopping/deploys)

> A modern, scalable microservices-based e-commerce API built with Go, featuring comprehensive testing, CI/CD, and beautiful API documentation.

### Live Demo - Deployed on Netlify
**[GoCart E-commerce Platform](https://gocartshopping.netlify.app/)** - Experience the full-featured shopping platform with real-time product browsing, user management, and order processing!

## **Live API Documentation**

**[Interactive API Explorer](https://wasifsarwar.github.io/gocart/)** - Test APIs directly in your browser

### Service-Specific Documentation
- **Product API**: [Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/wasifsarwar/gocart/main/api/product/openapi.yaml) | [OpenAPI Spec](api/product/openapi.yaml)
- **User API**: [Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/wasifsarwar/gocart/main/api/user/openapi.yaml) | [OpenAPI Spec](api/user/openapi.yaml)
- **Order API**: [Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/wasifsarwar/gocart/main/api/order-management/openapi.yaml) | [OpenAPI Spec](api/order-management/openapi.yaml)

---

## **Architecture Overview**

GoCart implements a **microservices architecture** with the following services:

### **Frontend Stack**
- **React 18** with TypeScript
- **React Router** for navigation  
- **Custom hooks** for state management
- **Responsive CSS** with mobile-first design
- **Real-time search & sorting**
- **Animated gopher branding**

### **Implemented Services**

| Service | Endpoint | Status | Description |
|---------|----------|--------|-------------|
| **Product Service** | `/products` | Complete | Product catalog, inventory, CRUD operations 
| **User Service**    | `/users`    | Complete | User management, authentication, profiles 
| **Order Service**   | `/orders`   | Complete | Order processing and management 

**All services run on single port**

### **Planned Services**
- **Payment Service** - Payment processing simulation
- **Notification Service** - Email/SMS notifications

---

## **Quick Start**

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

#### **Option A: Run All Services (Single Port)**
```bash
go run cmd/main.go
```
**Unified API**: http://localhost:8080
- Products: http://localhost:8080/products
- Users: http://localhost:8080/users  
- Orders: http://localhost:8080/orders

#### **Option B: Docker Compose**
```bash
# Start all services
docker-compose up

# Build and start
docker-compose up --build

# Stop all services
docker-compose down

# View logs
docker-compose logs app      # Backend logs
```

#### **Option C: Bash startup script**
```bash
# Start everything with one command
./start.sh
```


---

## **API Endpoints**

### **Product Service**
```http
GET    /products           # List all products
POST   /products           # Create product
GET    /products/{id}      # Get product by ID
PUT    /products/{id}      # Update product
DELETE /products/{id}      # Delete product
```

### **User Service** 
```http
GET    /users              # List all users
POST   /users/register     # Register new user
POST   /users/login        # Login user
GET    /users/{id}         # Get user by ID
PUT    /users/{id}         # Update user
DELETE /users/{id}         # Delete user
```

### **Order Service**
```http
GET    /orders             # List all orders
POST   /orders             # Create new order
GET    /orders/{id}        # Get order by ID
PUT    /orders/{id}        # Update order
DELETE /orders/{id}        # Delete order
GET    /orders/user/{user_id} # Get orders by user ID
DELETE /orders/{id}/items  # Delete order item
```


## **Testing**

### **Run All Tests**
```bash
go test ./...
```
