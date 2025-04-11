package middleware

import(
	"context"
	"net/http"
	"fmt"
	"strings"
	"github.com/shayantrix/task_manager_api/pkg/tokens"
	//"github.com/golang-jwt/jwt/v5"
)

// Define an authentication middleware that protects routes that need authentication
//next http.Handler
func Authentication(next http.Handler) http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		tokenString := r.Header.Get("Authorization")
		if tokenString == ""{
			fmt.Println("Error in token string reading")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Missing authentication token"}`))
			return 
		}

		// The token should be prefixed with "Bearer "
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts [0] != "Bearer" {
			fmt.Println("Error in Bearer part")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Missing authentication token"}`))
			return
		}

		tokenString = tokenParts[1]

		// Have the claims if the token is valid
		claims, err := tokens.ValidateJWT(tokenString)
		fmt.Println(claims)
		if err != nil{
			fmt.Println("Error in claims part")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Missing authentication token"}`))
			return
		}
		// If the token is valid, the user’s ID is stored in the Gin context for further use.
		// Add claims to the request context
		//userIDVal := claims.ID
		//fmt.Println(claims.ID)
		//if !exists{
		//	http.Error(w, "Unauthorized - user_id not in the token", http.StatusUnauthorized)
		//	return
		//}

		//userID, ok := userIDVal.(string)
		//if !ok {
		//	http.Error(w, "Unauthorized - user_id not a string", http.StatusUnauthorized)
		//	return
		//}
		ctx := context.WithValue(r.Context(), "id", claims.ID)
		//next.ServeHTTP(w, r.WithContext(ctx))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

