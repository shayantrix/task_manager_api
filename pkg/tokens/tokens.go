package tokens

import(
	"fmt"
	"log"
	"time"
        "github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("In the darkest night look for the light")


//type CustomClaims struct{
//	ID uuid.UUID `json:"id"`
	//expireTime  MyExpireTime `json:"expire_time"`
//}

// type Myfloat float64
//type MyExpireTime time.Now().Add(time.Hour*24).Unix()

func JWTGenerate(ID uuid.UUID) (string, error){
        // Make jwt object token
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
                "id": ID,
                "exp": time.Now().Add(time.Hour*24).Unix(),
        })

        jwt_string, err := token.SignedString(jwtKey)
        if err != nil{
                return "", err
        }
        return jwt_string, err
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
        //Parse The token
        // jwt.ParseWithClaims(token_String, &MycustomClaims, 
        //var myCustomClaims CustomClaims
	// Instead of &myCustomClaims we can use (jwt.MapClaims) if we didnt have custom tokens. The default is jwt.MapClaims
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
                //Check the login Method (jwt.SigningMethodHS256)
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("Error in siging method")
                }

                return jwtKey, nil
        })

        if err != nil{
                log.Fatal("Error in parsing the token: %s", err)
                return nil, fmt.Errorf("Error in Parsing the token: %s", err)
        }

        //Validate the token
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
                return claims, nil
        }

        return nil, fmt.Errorf("Invalid Token!")
}


