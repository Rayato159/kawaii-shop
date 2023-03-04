package repositories

import (
	"fmt"

	"github.com/Rayato159/kawaii-shop/modules/users"
	"github.com/Rayato159/kawaii-shop/modules/users/repositories/patterns"

	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetProfile(userId string) (*users.User, error)
}

type usersRepository struct {
	Db *sqlx.DB
}

func UsersRepository(db *sqlx.DB) IUsersRepository {
	return &usersRepository{Db: db}
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
