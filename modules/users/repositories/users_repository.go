package repositories

import (
	"context"
	"fmt"

	"github.com/Rayato159/kawaii-shop/modules/users"
	"github.com/Rayato159/kawaii-shop/modules/users/repositories/patterns"

	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	GetTransaction() (*sqlx.Tx, error)
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetProfile(userId string) (*users.User, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOauth(req *users.UserPassport) error
	DeleteOauth(refreshToken string) error
}

type usersRepository struct {
	Db *sqlx.DB
}

func UsersRepository(db *sqlx.DB) IUsersRepository {
	return &usersRepository{Db: db}
}

func (r *usersRepository) GetTransaction() (*sqlx.Tx, error) {
	tx, err := r.Db.BeginTxx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *usersRepository) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// Inserting
	result, err := patterns.InsertUser(r.Db, req).Customer()
	if err != nil {
		return nil, err
	}

	// Get result from inserting
	user, err := result.Result()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *usersRepository) GetProfile(userId string) (*users.User, error) {
	query := `
	SELECT
		"u"."id",
		"u"."email",
		"u"."username",
		"r"."title" AS "role"
	FROM "users" "u"
		LEFT JOIN "roles" "r" ON "r"."id" = "u"."role_id"
	WHERE "u"."id" = $1;`

	profile := new(users.User)
	if err := r.Db.Get(profile, query, userId); err != nil {
		return nil, fmt.Errorf("get user profile failed: %v", err)
	}
	return profile, nil
}

func (r *usersRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
	SELECT
		"u"."id",
		"u"."email",
		"u"."password",
		"u"."username",
		"r"."title" AS "role"
	FROM "users" "u"
		LEFT JOIN "roles" "r" ON "r"."id" = "u"."role_id"
	WHERE "u"."email" = $1;`

	user := new(users.UserCredentialCheck)
	if err := r.Db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *usersRepository) InsertOauth(req *users.UserPassport) error {
	query := `
	INSERT INTO "oauth" (
		"user_id",
		"refresh_token",
		"access_token"
	)
	VALUES ($1, $2, $3)
		RETURNING "id" AS "code";`

	if err := r.Db.QueryRowxContext(
		context.Background(),
		query,
		req.User.Id,
		req.Token.RefreshToken,
		req.Token.AccessToken,
	).Scan(&req.Token.Code); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) DeleteOauth(refreshToken string) error {
	query := `
	DELETE FROM "oauth"
		WHERE "refresh_token" = $1;`

	if _, err := r.Db.ExecContext(context.Background(), query, refreshToken); err != nil {
		return fmt.Errorf("refresh_token not found")
	}
	return nil
}
