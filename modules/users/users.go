package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserCredential struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type UserCredentialCheck struct {
	Id       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Username string `db:"username"`
	Role     string `db:"role"`
}

type UserRemoveCredential struct {
	Code string `json:"code" form:"code"`
}

type UserRefreshCredential struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type Oauth struct {
	Id     string `db:"id" json:"id" form:"id"`
	UserId string `db:"user_id"`
}

type UserPassport struct {
	User  *User      `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id           string `db:"id" json:"id"`
	AccessToken  string `db:"access_token" json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}

type User struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	Role     string `db:"role" json:"role"`
}

type UserClaims struct {
	Id string `db:"id" json:"id"`
}

type UserRegisterReq struct {
	Email    string `db:"email" json:"email" form:"email"`
	Username string `db:"username" json:"username" form:"username"`
	Password string `db:"password" json:"password" form:"password"`
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
