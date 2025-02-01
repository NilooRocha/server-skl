package middleware

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"server/domain"
)

type UserClaims struct {
	UserID struct {
		Subject string
	} `json:"id"`
	jwt.RegisteredClaims
}

func RequireAuth(next http.Handler, userRepo domain.IUser) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Error(w, "Authorization cookie is missing", http.StatusUnauthorized)
				return
			}
			http.Error(w, fmt.Sprintf("Error retrieving cookie: %v", err), http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// TODO: move secret key to env
			return []byte("5asfg67sdftgs57df4g5764sdfg473sd4f62g6sdf3sd2g46sdf352sdf4"), nil
		})

		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing token: %v", err), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		user, err := userRepo.Read(claims.Subject)
		if err != nil {
			http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusUnauthorized)
			return
		}

		if !user.IsVerified {
			//TODO: send verification email
			http.Error(w, "Email is not verified", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
