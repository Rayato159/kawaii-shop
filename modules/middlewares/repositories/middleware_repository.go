package repositories

import "github.com/jmoiron/sqlx"

type IMiddlewareRepository interface {
	FindAccessToken(userId string, accessToken string) bool
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
