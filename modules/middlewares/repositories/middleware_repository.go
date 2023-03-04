package repositories

import (
	"fmt"

	"github.com/Rayato159/kawaii-shop/modules/middlewares"

	"github.com/jmoiron/sqlx"
)

type IMiddlewareRepository interface {
	FindAccessToken(userId string, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewareRepository struct {
	Db *sqlx.DB
}

func MiddlewareRepository(db *sqlx.DB) IMiddlewareRepository {
	return &middlewareRepository{
		Db: db,
	}
}

func (r *middlewareRepository) FindAccessToken(userId string, accessToken string) bool {
	query := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "oauth"
	WHERE "user_id" = $1
	AND "access_token" = $2;`

	var check bool
	if err := r.Db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}
	return check
}

func (r *middlewareRepository) FindRole() ([]*middlewares.Role, error) {
	query := `
	SELECT
		"id",
		"title"
	FROM "roles"
	ORDER BY "id" DESC;`

	roles := make([]*middlewares.Role, 0)
	if err := r.Db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are empty")
	}
	return roles, nil
}
