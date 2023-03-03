package repositories

import (
	"fmt"

	"github.com/Rayato159/kawaii-shop/modules/oauth"
	"github.com/Rayato159/kawaii-shop/modules/oauth/repositories/patterns"

	"github.com/jmoiron/sqlx"
)

type IOauthRepository interface {
	InsertCustomer(req *oauth.UserRegisterReq) (*oauth.UserPassport, error)
	GetProfile(userId string) (*oauth.User, error)
}

type oauthRepository struct {
	Db *sqlx.DB
}

func OauthRepository(db *sqlx.DB) IOauthRepository {
	return &oauthRepository{Db: db}
}

func (r *oauthRepository) InsertCustomer(req *oauth.UserRegisterReq) (*oauth.UserPassport, error) {
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

func (r *oauthRepository) GetProfile(userId string) (*oauth.User, error) {
	query := `
	SELECT
		"u"."id",
		"u"."email",
		"u"."username",
		"r"."title" AS "role"
	FROM "users" "u"
		LEFT JOIN "roles" "r" ON "r"."id" = "u"."role_id"
	WHERE "u"."id" = $1;`

	profile := new(oauth.User)
	if err := r.Db.Get(profile, query, userId); err != nil {
		return nil, fmt.Errorf("get user profile failed: %v", err)
	}
	return profile, nil
}
