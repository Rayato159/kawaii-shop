package usecases

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Rayato159/kawaii-shop/config"
	filespkg "github.com/Rayato159/kawaii-shop/modules/files"
)

type IFilesUsecase interface {
	UploadToGCP(req []*filespkg.FileReq) ([]*filespkg.FileRes, error)
}

type filesUsecase struct {
	cfg config.IConfig
}

func FilesUsecase(cfg config.IConfig) IFilesUsecase {
	return &filesUsecase{
		cfg: cfg,
	}
}

func (u *filesUsecase) UploadToGCP(req []*filespkg.FileReq) ([]*filespkg.FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	res := make([]*filespkg.FileRes, 0)
	for i := range req {
		container, err := req[i].File.Open()
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(container)
		if err != nil {
			return nil, err
		}

		buf := bytes.NewBuffer(b)

		// Upload an object with storage.Writer.
		wc := client.Bucket(u.cfg.App().GCPBucket()).Object(req[i].Destination).NewWriter(ctx)

		if _, err = io.Copy(wc, buf); err != nil {
			return nil, fmt.Errorf("io.Copy: %v", err)
		}
		// Data can continue to be added to the file until the writer is closed.
		if err := wc.Close(); err != nil {
			return nil, fmt.Errorf("Writer.Close: %v", err)
		}
		log.Printf("%v uploaded to %v.\n", req[i].FileName, req[i].Destination)

		res = append(res, &filespkg.FileRes{
			Url:      fmt.Sprintf("https://storage.googleapis.com/%s/%s", u.cfg.App().GCPBucket(), req[i].Destination),
			Filename: req[i].FileName,
		})
	}
	return res, nil
}
