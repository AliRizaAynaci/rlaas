package server

import (
	"encoding/json"
	"github.com/AliRizaAynaci/rlaas/internal/middleware"
	"github.com/AliRizaAynaci/rlaas/internal/server/handlers"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/health", s.HealthHandler)

	// --- Public Auth & Check ---
	mux.HandleFunc("/auth/google/login", handlers.LoginHandler)       // GET  /auth/google/login
	mux.HandleFunc("/auth/google/callback", handlers.CallbackHandler) // GET  /auth/google/callback
	mux.HandleFunc("/check", handlers.RateLimitCheck)                 // POST /check
	mux.HandleFunc("/logout", handlers.LogoutHandler)                 // GET  /logout

	// --- Protected Auth endpoints ---
	// Get current user info (requires authentication)
	mux.Handle("/me",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.MeHandler),
		),
	)

	// --- Project Creation ---
	// Creates a new project and returns API key
	mux.Handle("/register",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.RegisterProjectHandler),
		),
	)

	// --- Rate-Limit Rules endpoints requires an API-KEY ---
	// 1. List all rules for the API-key in Authorization header
	mux.Handle("/rules",
		http.HandlerFunc(handlers.GetRules),
	)

	// 2. Add a new rule
	mux.Handle("/rule/add",
		http.HandlerFunc(handlers.AddRule),
	)

	// 3. Update or delete an existing rule by ID:
	//    PUT    /rule/{id}
	//    DELETE /rule/{id}
	mux.Handle("/rule/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPut:
				handlers.UpdateRule(w, r) // handle update
			case http.MethodDelete:
				handlers.DeleteRule(w, r) // handle delete
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		}),
	)

	mux.Handle("/projects",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.ListProjectsHandler),
		),
	)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for frontend integration
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Frontend URL
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Enable credentials for cookie-based auth

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
