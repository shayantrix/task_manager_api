package auth

import(
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error){
	password_hashed , err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(password_hashed), err
}

func CheckHashedPassword(password, hash_password string) error{
	if err := bcrypt.CompareHashAndPassword([]byte(hash_password), []byte(password)); err != nil{
		return err
	}
	return nil
}

