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

	// --- Project Creation ---
	// Creates a new project and returns API key
	mux.Handle("/register",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.RegisterProjectHandler),
		),
	)

	// --- Rate-Limit Rules (aka Endpoints) ---
	// List all rules for the API-key in Authorization header
	mux.Handle("/rules",
		middleware.AuthMiddleware(http.HandlerFunc(handlers.GetRules)),
	)
	// Add a new rule
	mux.Handle("/rule/add",
		middleware.AuthMiddleware(http.HandlerFunc(handlers.AddRule)),
	)
	// Update or delete an existing rule by ID:
	//   PUT    /rule/{id}
	//   DELETE /rule/{id}
	mux.Handle("/rule/",
		middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPut:
				handlers.UpdateRule(w, r)
			case http.MethodDelete:
				handlers.DeleteRule(w, r)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})),
	)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

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
