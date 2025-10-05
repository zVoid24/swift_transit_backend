package user

import (
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents JSON response
type LoginResponse struct {
	Message string      `json:"message"`
	User    interface{} `json:"user,omitempty"`
}

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Only allow POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(LoginResponse{Message: "Method not allowed"})
		return
	}

	// Decode JSON request
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{Message: "Invalid JSON payload"})
		return
	}

	// Validate input
	if req.UserName == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{Message: "Username and password required"})
		return
	}

	// Check credentials
	user, err := h.UserRepo.Find(req.UserName, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{Message: "Invalid username or password"})
		return
	}

	// Success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Message: "Login successful",
		User:    user,
	})
}
