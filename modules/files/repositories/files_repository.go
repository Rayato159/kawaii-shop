package repositories

import "github.com/jmoiron/sqlx"

type IFilesRepository interface{}

type fileRepository struct {
	db *sqlx.DB
}

func FilesRepository(db *sqlx.DB) IFilesRepository {
	return &fileRepository{db: db}
}
