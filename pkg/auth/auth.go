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

/*
import(
	"github.com/shayantrix/task_manager_api/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type SecureAuth struct{
	Email	string	`json:"email"`
	Pass  []byte	`json:"-"`
}

var secureAuth []SecureAuth

func EvaluatePass(reg *models.RegisterData)error{
	//reg == ID, Name, Email, Pass
	for _, v := range models.registerData{
		secureAuth = append(secureAuth, SecureAuth{
			Email: v.Email,
			Pass: bcrypt.GenerateFromPassword(v.Pass, bcrypt.DefaultCost),
		})
	}


	
	//func CompareHashAndPassword(hashedPassword, password []byte) error

	for _, v := range secureAuth{
		if err := bcrypt.CompareHashAndPassword(v.Pass, reg.Pass); err != nil{
			return err
		}
	}
	return nil

}
*/


