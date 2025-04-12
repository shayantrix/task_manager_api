package tokens

import(
	"os"
	"fmt"
	"time"
        "github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)
//define it in env
var jwtKey = []byte(os.Getenv("JWT_TOKEN"))


type CustomClaims struct{
	ID uuid.UUID `json:"id"`
	//expireTime  MyExpireTime `json:"expire_time"`
	jwt.RegisteredClaims
}

// type Myfloat float64
//type MyExpireTime time.Now().Add(time.Hour*24).Unix()

func JWTGenerate(IDl uuid.UUID) (string, error){
	// For custom claims 
	claims := CustomClaims {
		IDl, 
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("Error in JWT token generating: %v", err)
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (*CustomClaims, error) {
        //Parse The token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error){
                //Check the login Method (jwt.SigningMethodHS256)
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("Error in siging method")
                }

                return jwtKey, nil
        })

        if err != nil || !token.Valid{
                return nil, fmt.Errorf("Error in Parsing the token: %s", err)
        }

	if claims, ok := token.Claims.(*CustomClaims); ok {
		fmt.Println(claims.ID)
		return claims, nil
	}else{
		return nil, fmt.Errorf("Error in recieving claims: %v", err)
	}
}


