package entities

type UserCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	Token    string `db:"token" json:"token"`
	Role     string `db:"role" json:"role"`
}

type UserRegisterReq struct {
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Token    string `db:"token" json:"token"`
}

type AdminRegisterReq struct {
	Email      string `db:"email" json:"email"`
	Username   string `db:"username" json:"username"`
	Password   string `db:"password" json:"password"`
	Token      string `db:"token" json:"token"`
	AdminToken string `db:"admin_token" json:"admin_token"`
}
