package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

type Claims struct {
	PasswordHash string `json:"pwd_hash"`
	jwt.RegisteredClaims
}

func (a *Api) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, AuthResponse{Error: "Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, AuthResponse{Error: "Invalid request"}, http.StatusBadRequest)
		return
	}

	envPassword := os.Getenv("TODO_PASSWORD")
	if envPassword == "" {
		writeJSON(w, AuthResponse{Error: "Auth password is not set"}, http.StatusInternalServerError)
		return
	}

	if req.Password != envPassword {
		writeJSON(w, AuthResponse{Error: "Invalid password"}, http.StatusUnauthorized)
		return
	}

	hash := sha256.Sum256([]byte(envPassword))
	passwordHash := hex.EncodeToString(hash[:])

	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &Claims{
		PasswordHash: passwordHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(envPassword))
	if err != nil {
		writeJSON(w, AuthResponse{Error: "Failed to generate token"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, AuthResponse{Token: tokenString}, http.StatusOK)
}

func (a *Api) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		envPassword := os.Getenv("TODO_PASSWORD")
		if envPassword == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			writeJSON(w, AuthResponse{Error: "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(envPassword), nil
		})

		if err != nil || !token.Valid {
			writeJSON(w, AuthResponse{Error: "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		hash := sha256.Sum256([]byte(envPassword))
		currentHash := hex.EncodeToString(hash[:])
		if claims.PasswordHash != currentHash {
			writeJSON(w, AuthResponse{Error: "Unauthorized"}, http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
