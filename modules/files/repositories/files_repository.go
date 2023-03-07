package repositories

import (
	filespkg "github.com/Rayato159/kawaii-shop/modules/files"

	"github.com/jmoiron/sqlx"
)

type IFilesRepository interface{}

type fileRepository struct {
	db *sqlx.DB
}

func FilesRepository(db *sqlx.DB) IFilesRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) UploadToGCP(req *filespkg.FileReq) (*filespkg.FileRes, error) {
	return nil, nil
}
