# GoCart E-commerce Microservice API

[![API Documentation](https://img.shields.io/badge/API-Documentation-blue)](https://wasifsarwar.github.io/gocart/)
[![Product API](https://img.shields.io/badge/Product%20API-Swagger-green)](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/wasifsarwar/gocart/main/api/product/openapi.yaml)
[![User API](https://img.shields.io/badge/User%20API-Swagger-green)](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/wasifsarwar/gocart/main/api/user/openapi.yaml)

## 📚 API Documentation

**🔗 [Live API Documentation](https://wasifsarwar.github.io/gocart/)** - Interactive Swagger UI

### Individual Service APIs
- **Product Service**: [Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/wasifsarwar/gocart/main/api/product/openapi.yaml) | [Raw Spec](api/product/openapi.yaml)
- **User Service**: [Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/wasifsarwar/gocart/main/api/user/openapi.yaml) | [Raw Spec](api/user/openapi.yaml)

## Project Overview
Build a microservice-based e-commerce API system with the following components:

Product Service

Product catalog management
Inventory tracking
Category management
Search functionality


User Service

Authentication (JWT implementation)
User management
Role-based access control


Order Service

Order processing
Payment integration simulation
Order history and tracking



Key Technical Requirements

Microservice Architecture: Implement separate services that communicate via REST or gRPC
Database Integration: Use PostgreSQL and implement proper migrations
Caching Layer: Implement Redis for caching product information
Docker Integration: Create Docker files and Docker Compose for easy deployment
Testing: Comprehensive unit and integration tests
Documentation: Well-documented API with Swagger
Logging & Monitoring: Implement structured logging and basic metrics
Rate Limiting: Implement rate limiting for API endpoints

Stretch Goals

Circuit breaker pattern
Message queue integration (RabbitMQ or Kafka)
Basic CI/CD pipeline configuration


![project structure](/docs/images/Project%20structure.png)
```
/gocart
└── product-service
    ├── cmd
    │   └── main.go                # Entry point for the service
    ├── internal
    │   ├── models
    │   │   └── product.go         # Product model definition
    │   ├── handlers
    │   │   └── product_handler.go  # HTTP handlers for product-related endpoints
    │   ├── repository
    │   │   └── product_repository.go # Database interactions for products
    │   ├── services
    │   │   └── product_service.go  # Business logic for product operations
    │   └── middleware
    │       └── auth_middleware.go  # Middleware for authentication/authorization
    ├── config
    │   └── config.go              # Configuration management (e.g., loading environment variables)
    ├── Dockerfile                  # Dockerfile for building the service
    ├── docker-compose.yml          # (Optional) If you want to run the service with other services
    ├── go.mod                      # Go module file for dependency management
    ├── go.sum                      # Go module checksum file
    └── README.md                   # Documentation for the product service
```


```
sql commands
admin: psql -U wasifsmacbookpro -h localhost -p 5432 -d postgres

CREATE TABLE products (
    id SERIAL PRIMARY KEY,          -- Unique identifier for each product
    name VARCHAR(255) NOT NULL,     -- Name of the product
    description TEXT,                -- Description of the product
    price NUMERIC(10, 2) NOT NULL,   -- Price of the product (up to 10 digits, 2 decimal places)
    category VARCHAR(100),           -- Category of the product
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp for when the product was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Timestamp for when the product was last updated
);
```

