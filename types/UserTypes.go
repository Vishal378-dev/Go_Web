package types

import (
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID       interface{} `bson:"_id,omitempty"`
	Name     string      `json:"name" bson:"name"`
	Email    string      `json:"email" bson:"email"`
	Phone    string      `json:"phone" bson:"phone"`
	Password string      `json:"password" bson:"password"`
	Role     string      `json:"role" bson:"role"`
	Created  time.Time   `json:"created_at" bson:"created_at"`
	Updated  time.Time   `json:"updated_at" bson:"updated_at"`
}

// for jwt
type UserClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func (u *User) ValidateRequest() error {
	regex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	validEmailReg := regexp.MustCompile(regex)
	if u.Email == "" || !validEmailReg.MatchString(u.Email) {
		return fmt.Errorf("invalid Email")
	}
	if u.Name == "" || len(u.Name) < 3 {
		return fmt.Errorf("invalid UserName")
	}
	if u.Password == "" || len(u.Password) < 7 {
		return fmt.Errorf("invalid Password")
	}
	if u.Phone == "" || len(u.Phone) != 10 {
		return fmt.Errorf("invalid Phone")
	}
	return nil
}

// NewUser is a constructor function that returns a User with default values
func NewUser() User {
	return User{
		Name:     "",
		Email:    "",
		Phone:    "",
		Password: "",
		Role:     "user",
		Created:  time.Now(),
		Updated:  time.Now(),
	}
}
