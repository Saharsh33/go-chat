package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"chat-server/internal/auth"
	postgres "chat-server/internal/storage/db"
)

// StorageAuth is the subset of the storage interface needed by the auth handlers.
type StorageAuth interface {
	RegisterUser(ctx context.Context, username, passwordHash string) error
	GetUserPasswordHash(ctx context.Context, username string) (string, error)
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// NewRegisterHandler returns an http.HandlerFunc for POST /register.
func NewRegisterHandler(store StorageAuth, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req authRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
			return
		}

		if req.Username == "" || req.Password == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "username and password are required"})
			return
		}

		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "could not hash password"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := store.RegisterUser(ctx, req.Username, hash); err != nil {
			if err == postgres.ErrUserExists {
				writeJSON(w, http.StatusConflict, errorResponse{Error: "username already taken"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "could not register user"})
			return
		}

		token, err := auth.GenerateToken(req.Username, jwtSecret)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "could not generate token"})
			return
		}

		writeJSON(w, http.StatusCreated, authResponse{Token: token})
	}
}

// NewLoginHandler returns an http.HandlerFunc for POST /login.
func NewLoginHandler(store StorageAuth, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req authRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
			return
		}

		if req.Username == "" || req.Password == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "username and password are required"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		hash, err := store.GetUserPasswordHash(ctx, req.Username)

		// Always run the bcrypt comparison to prevent timing-based enumeration of
		// existing usernames, even when the user does not exist.
		hashToCheck := hash
		if err != nil || hash == "" {
			hashToCheck = auth.DummyHash
		}

		valid := auth.CheckPassword(hashToCheck, req.Password)
		if err != nil || hash == "" || !valid {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid credentials"})
			return
		}

		token, err := auth.GenerateToken(req.Username, jwtSecret)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "could not generate token"})
			return
		}

		writeJSON(w, http.StatusOK, authResponse{Token: token})
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}
