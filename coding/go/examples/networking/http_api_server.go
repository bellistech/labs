// HTTP API Server - A simple REST API example
//
// This example builds a JSON API server demonstrating:
// - HTTP routing with net/http
// - JSON encoding/decoding
// - Middleware pattern
// - Error handling
// - Request context
// - Graceful shutdown
//
// Usage:
//   go run http_api_server.go
//
// Test endpoints:
//   curl http://localhost:8080/health
//   curl http://localhost:8080/api/users
//   curl -X POST -d '{"name":"Alice","email":"alice@example.com"}' http://localhost:8080/api/users
//   curl http://localhost:8080/api/users/1
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ============================================================
// Models
// ============================================================

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// ============================================================
// In-memory store (would be a database in production)
// ============================================================

type UserStore struct {
	mu     sync.RWMutex
	users  map[int]*User
	nextID int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (s *UserStore) Create(name, email string) *User {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	user := &User{
		ID:        s.nextID,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
	s.users[user.ID] = user
	s.nextID++
	return user
}

func (s *UserStore) Get(id int) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[id]
	return user, ok
}

func (s *UserStore) List() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	users := make([]*User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}

func (s *UserStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, ok := s.users[id]; !ok {
		return false
	}
	delete(s.users, id)
	return true
}

// ============================================================
// API Server
// ============================================================

type APIServer struct {
	store  *UserStore
	router *http.ServeMux
}

func NewAPIServer() *APIServer {
	s := &APIServer{
		store:  NewUserStore(),
		router: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *APIServer) routes() {
	// Health check
	s.router.HandleFunc("/health", s.handleHealth)
	
	// API routes
	s.router.HandleFunc("/api/users", s.handleUsers)
	s.router.HandleFunc("/api/users/", s.handleUser)
}

// ServeHTTP implements http.Handler
func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Wrap with middleware
	handler := s.loggingMiddleware(s.router)
	handler.ServeHTTP(w, r)
}

// ============================================================
// Middleware
// ============================================================

func (s *APIServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create response wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, status: 200}
		
		next.ServeHTTP(wrapped, r)
		
		log.Printf("%s %s %d %v",
			r.Method, r.URL.Path, wrapped.status, time.Since(start))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// ============================================================
// Handlers
// ============================================================

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.methodNotAllowed(w)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *APIServer) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listUsers(w, r)
	case http.MethodPost:
		s.createUser(w, r)
	default:
		s.methodNotAllowed(w)
	}
}

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path: /api/users/{id}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	
	switch r.Method {
	case http.MethodGet:
		s.getUser(w, r, id)
	case http.MethodDelete:
		s.deleteUser(w, r, id)
	default:
		s.methodNotAllowed(w)
	}
}

func (s *APIServer) listUsers(w http.ResponseWriter, r *http.Request) {
	users := s.store.List()
	s.jsonResponse(w, http.StatusOK, users)
}

func (s *APIServer) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		s.jsonError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	
	if input.Name == "" || input.Email == "" {
		s.jsonError(w, http.StatusBadRequest, "name and email required")
		return
	}
	
	user := s.store.Create(input.Name, input.Email)
	s.jsonResponse(w, http.StatusCreated, user)
}

func (s *APIServer) getUser(w http.ResponseWriter, r *http.Request, id int) {
	user, ok := s.store.Get(id)
	if !ok {
		s.jsonError(w, http.StatusNotFound, "user not found")
		return
	}
	s.jsonResponse(w, http.StatusOK, user)
}

func (s *APIServer) deleteUser(w http.ResponseWriter, r *http.Request, id int) {
	if !s.store.Delete(id) {
		s.jsonError(w, http.StatusNotFound, "user not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ============================================================
// Response helpers
// ============================================================

func (s *APIServer) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *APIServer) jsonError(w http.ResponseWriter, status int, message string) {
	s.jsonResponse(w, status, ErrorResponse{
		Error: message,
		Code:  status,
	})
}

func (s *APIServer) methodNotAllowed(w http.ResponseWriter) {
	s.jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
}

// ============================================================
// Main
// ============================================================

func main() {
	// Create server
	api := NewAPIServer()
	
	// Seed with some data
	api.store.Create("Bob", "bob@example.com")
	api.store.Create("Carol", "carol@example.com")
	
	// Create HTTP server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      api,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start server in background
	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Print usage
	fmt.Println()
	fmt.Println("API Endpoints:")
	fmt.Println("  GET    /health           - Health check")
	fmt.Println("  GET    /api/users        - List all users")
	fmt.Println("  POST   /api/users        - Create user (JSON body)")
	fmt.Println("  GET    /api/users/{id}   - Get user by ID")
	fmt.Println("  DELETE /api/users/{id}   - Delete user")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  curl http://localhost:8080/health")
	fmt.Println("  curl http://localhost:8080/api/users")
	fmt.Println("  curl -X POST -H 'Content-Type: application/json' \\")
	fmt.Println("       -d '{\"name\":\"Alice\",\"email\":\"alice@example.com\"}' \\")
	fmt.Println("       http://localhost:8080/api/users")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()
	
	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	
	log.Println("Shutting down...")
	
	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
	
	log.Println("Server stopped")
}
