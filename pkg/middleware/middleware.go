package middleware

import (	
	//"context"
	"github.com/shayantrix/task_manager_api/pkg/models"
	"net/http"
	"fmt"
	//"github.com/gorilla/mux"
	//"fmt"
	"strings"
	//"strconv"
	"github.com/gin-gonic/gin"
	"github.com/shayantrix/task_manager_api/pkg/tokens"
	//"github.com/golang-jwt/jwt/v5"
	//"github.com/shayantrix/task_manager_api/pkg/controllers"
)

// Define an authentication middleware that protects routes that need authentication
//next http.Handler
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context){
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
         		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Missing authentication token"})
            		c.Abort()
            		return
        	}
		// The token should be prefixed with "Bearer "
		tokenParts := strings.Split(tokenString, " ")
                if len(tokenParts) != 2 || tokenParts [0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Missing authentication token"})
			c.Abort()
			return
                }

		tokenString = tokenParts[1]

                // Have the claims if the token is valid
                claims, err := tokens.ValidateJWT(tokenString)
                if err != nil{
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid Authentication token"})
                        c.Abort()
                        return
                }
		c.Set("id", claims.ID)
		c.Next()
	}
}


/*
func Authentication(next http.Handler) http.HandlerFunc{
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
		ctx := context.WithValue(r.Context(), "id", claims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
*/

func Authorization() gin.HandlerFunc{
	return func(c *gin.Context){
                tokenString := c.GetHeader("Authorization")
                if tokenString == "" {
                        c.JSON(http.StatusUnauthorized, gin.H{"Error": "Missing authentication token"})
                        c.Abort()
                        return
                }
                // The token should be prefixed with "Bearer "
                tokenParts := strings.Split(tokenString, " ")
                if len(tokenParts) != 2 || tokenParts [0] != "Bearer" {
                        c.JSON(http.StatusUnauthorized, gin.H{"Error": "Missing authentication token"})
                        c.Abort()
                        return
                }

                tokenString = tokenParts[1]

                // Have the claims if the token is valid
                claims, err := tokens.ValidateJWT(tokenString)
                if err != nil{
                        c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid Authentication token"})
                        c.Abort()
                        return
                }
		found := false

                AllUsers := models.GetAllUsers()
                for _, item := range AllUsers{
                        if item.ID == claims.ID{
                                found = true
                                break
                        }
                }

		if !found{
                	c.JSON(http.StatusUnauthorized, gin.H{"Error": "Missing user's id"})
                        c.Abort()
                        return

		}
		
		c.Set("id", claims.ID)
		fmt.Println(claims.ID)
		c.Next()
	}
}


/*
func Authorization(next http.Handler) http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		tokenString := r.Header.Get("Authorization")
		if tokenString == ""{
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Missing authentication token"}`))
			return
		}
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Missing authentication token"}`))
			return
		}

		tokenString = tokenParts[1]

		claims, err := tokens.ValidateJWT(tokenString)
		if err != nil{
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Missing authentication token"}`))
			return 
		}
		/*	
		params := mux.Vars(r)
		paramID64, err := strconv.ParseInt(params["id"], 10, 0)
		paramID := int(paramID64)
		//Check if the userID from parameter is equal to JWT token id
		// Error in this part		
		//|| item.ID != claims.ID
		for i, item := range controllers.Data{
			if i != paramID || item.ID != claims.ID{
				w.WriteHeader(http.StatusForbidden)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"error": "Mismatching user id"}`))
				return 
			}
		}
		*/
/*
		found := false

		AllUsers := models.GetAllUsers()
		for _, item := range AllUsers{
			if item.ID == claims.ID{
				found = true
				break
			}
		}

		//for _, item := range controllers.Data{
		//	if item.ID == claims.ID{
		//		found = true
		//		break
		//	}
		//}

		if !found{
			w.WriteHeader(http.StatusForbidden)
                       	w.Header().Set("Content-Type", "application/json")
                        w.Write([]byte(`{"error": "Mismatching user id"}`))
                        return
		}

		ctx := context.WithValue(r.Context(), "id", claims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
*/
