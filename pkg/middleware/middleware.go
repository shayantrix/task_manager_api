package middleware

import(
	"context"
	"net/http"
	"strings"
	"github.com/shayantrix/task_manager_api/pkg/tokens"
	//"github.com/golang-jwt/jwt/v5"
)

// Define an authentication middleware that protects routes that need authentication
//next http.Handler
func Authentication() http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		tokenString := r.Header.Get("Authorization")
		if tokenString == ""{
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Missing authentication token"}`))
			return 
		}

		// The token should be prefixed with "Bearer "
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts [0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Missing authentication token"}`))
			return
		}

		tokenString = tokenParts[1]

		// Have the claims if the token is valid
		claims, err := tokens.ValidateJWT(tokenString)
		if err != nil{
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Missing authentication token"}`))
			return
		}
		// If the token is valid, the user’s ID is stored in the Gin context for further use.
		// Add claims to the request context
		ctx := context.WithValue(r.Context(), "userClaims", claims)
		//next.ServeHTTP(w, r.WithContext(ctx))
		var next http.Handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

