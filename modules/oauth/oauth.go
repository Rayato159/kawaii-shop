package oauth

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserCredential struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type UserPassport struct {
	User  *User      `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	Role     string `db:"role" json:"role"`
}

type UserRegisterReq struct {
	Email    string  `db:"email" json:"email" form:"email"`
	Username string  `db:"username" json:"username" form:"username"`
	Password string  `db:"password" json:"password" form:"password"`
	Token    *string `json:"token"`
}

func (obj *UserRegisterReq) BcryptHashing() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hashed password failed: %v", err)
	}
	obj.Password = string(hashedPassword)
	return nil
}

func (obj *UserRegisterReq) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.Email)
	if err != nil {
		return false
	}
	return match
}
