package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Serve OpenAPI YAML files
	r.PathPrefix("/api-specs/").Handler(http.StripPrefix("/api-specs/", http.FileServer(http.Dir("../api/"))))

	// API documentation endpoints
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>GoCart API Documentation</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .service { margin: 20px 0; padding: 20px; border: 1px solid #ddd; }
        a { color: #007bff; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <h1>GoCart API Documentation</h1>
    
    <div class="service">
        <h2>Product Service</h2>
        <p>Manages product catalog, inventory, and categories</p>
        <a href="https://petstore.swagger.io/?url=http://localhost:3000/api-specs/product/openapi.yaml">View in Swagger UI</a> |
        <a href="/api-specs/product/openapi.yaml">Raw OpenAPI Spec</a>
    </div>
    
    <div class="service">
        <h2>User Service</h2>
        <p>Handles user management, authentication, and profiles</p>
        <a href="https://petstore.swagger.io/?url=http://localhost:3000/api-specs/user/openapi.yaml">View in Swagger UI</a> |
        <a href="/api-specs/user/openapi.yaml">Raw OpenAPI Spec</a>
    </div>
    
    <div class="service">
        <h2>Development Servers</h2>
        <p>Product Service: <a href="http://localhost:8080">http://localhost:8080</a></p>
        <p>User Service: <a href="http://localhost:8081">http://localhost:8081</a></p>
    </div>
    
    <div class="service">
        <h2>Quick Start</h2>
        <p>1. Start your services: <code>go run cmd/main.go</code></p>
        <p>2. Visit the Swagger UI links above to test endpoints</p>
        <p>3. Import specs into Postman for testing</p>
    </div>
</body>
</html>
		`))
	})

	log.Println("Documentation server starting on :3000")
	log.Println("Visit http://localhost:3000 for API documentation")
	log.Fatal(http.ListenAndServe(":3000", r))
}
