package repositories

import (
	"github.com/Rayato159/kawaii-shop/modules/oauth"
	"github.com/Rayato159/kawaii-shop/modules/oauth/repositories/patterns"

	"github.com/jmoiron/sqlx"
)

type IOauthRepository interface {
	InsertCustomer(req *oauth.UserRegisterReq) (*oauth.UserPassport, error)
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
