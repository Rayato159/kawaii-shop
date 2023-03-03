package patterns

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Rayato159/kawaii-shop/modules/oauth"
	"github.com/jmoiron/sqlx"
)

type IInsertUser interface {
	tokenDecode() int
	Customer() (IInsertUser, error)
	Admin() (IInsertUser, error)
	Result() (*oauth.UserPassport, error)
}

func (f *userReq) tokenDecode() int {
	token := f.req.Token
	_ = token
	return -1
}

func (f *userReq) Customer() (IInsertUser, error) {
	ctx := context.Background()

	query := `
	INSERT INTO "users" (
		"email",
		"password",
		"username",
		"role_id"
	)
	VALUES
		($1, $2, $3, 1)
	RETURNING "id";`

	if err := f.db.QueryRowxContext(
		ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.Username,
	).Scan(&f.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username have been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email have been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}
	return f, nil
}

func (f *userReq) Admin() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO "users" (
		"email",
		"password",
		"username",
		"role_id"
	)
	VALUES
		($1, $2, $3, 2)
	RETURNING "id";`

	if err := f.db.QueryRowxContext(
		ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.Username,
	).Scan(&f.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username have been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email have been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}
	return f, nil
}

func (f *userReq) Result() (*oauth.UserPassport, error) {
	query := `
	SELECT
			json_build_object(
				'user', "t",
				'token', NULL
			) 
		FROM (
			SELECT
				"u"."id",
				"u"."email",
				"u"."username",
				"r"."title" AS "role"
			FROM "users" "u"
			LEFT JOIN "roles" "r" ON "r"."id" = "u"."role_id"
			WHERE "u"."id" = $1
	) AS "t";`

	data := make([]byte, 0)
	if err := f.db.Get(&data, query, f.id); err != nil {
		return nil, fmt.Errorf("get users failed: %v", err)
	}

	user := new(oauth.UserPassport)
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}
	return user, nil
}

type userReq struct {
	id  string
	req *oauth.UserRegisterReq
	db  *sqlx.DB
}

type customer struct {
	*userReq
}

type admin struct {
	*userReq
}

func InsertUser(db *sqlx.DB, req *oauth.UserRegisterReq) IInsertUser {
	token := req.Token
	_ = token
	roleId := 0

	switch roleId {
	case 2:
		return newAdmin(db, req)
	default:
		return newCustomer(db, req)
	}
}

func newAdmin(db *sqlx.DB, req *oauth.UserRegisterReq) IInsertUser {
	return &admin{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func newCustomer(db *sqlx.DB, req *oauth.UserRegisterReq) IInsertUser {
	return &customer{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}
