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
	DeleteFileInGCP(req []*filespkg.DeleteFileReq) error
}

type filesUsecase struct {
	cfg config.IConfig
}

type fileRes struct {
	bucket      string
	destination string
	file        *filespkg.FileRes
}

func FilesUsecase(cfg config.IConfig) IFilesUsecase {
	return &filesUsecase{
		cfg: cfg,
	}
}

func (f *fileRes) public() error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	acl := client.Bucket(f.bucket).Object(f.destination).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %v", err)
	}
	fmt.Printf("blob %v is now publicly accessible.\n", f.destination)
	return nil
}

func (u *filesUsecase) uploadWorkers(ctx context.Context, client *storage.Client, jobs <-chan *filespkg.FileReq, results chan<- *filespkg.FileRes, errChan chan<- error) {
	for job := range jobs {
		container, err := job.File.Open()
		if err != nil {
			errChan <- err
			return
		}
		b, err := ioutil.ReadAll(container)
		if err != nil {
			errChan <- err
			return
		}

		buf := bytes.NewBuffer(b)

		// Upload an object with storage.Writer.
		wc := client.Bucket(u.cfg.App().GCPBucket()).Object(job.Destination).NewWriter(ctx)

		if _, err = io.Copy(wc, buf); err != nil {
			errChan <- fmt.Errorf("io.Copy: %v", err)
			return
		}
		// Data can continue to be added to the file until the writer is closed.
		if err := wc.Close(); err != nil {
			errChan <- fmt.Errorf("Writer.Close: %v", err)
			return
		}
		log.Printf("%v uploaded to %v.\n", job.FileName, job.Destination)

		newFile := &fileRes{
			file: &filespkg.FileRes{
				Url:      fmt.Sprintf("https://storage.googleapis.com/%s/%s", u.cfg.App().GCPBucket(), job.Destination),
				Filename: job.FileName,
			},
			destination: job.Destination,
			bucket:      u.cfg.App().GCPBucket(),
		}

		// Make obj to public access
		if err := newFile.public(); err != nil {
			errChan <- err
			return
		}

		// Assign result
		errChan <- nil
		results <- newFile.file
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

	// Chan declare
	jobsCh := make(chan *filespkg.FileReq, len(req))
	resultsCh := make(chan *filespkg.FileRes, len(req))
	errsCh := make(chan error, len(req))

	// Return value
	res := make([]*filespkg.FileRes, 0)

	for _, f := range req {
		jobsCh <- f
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go u.uploadWorkers(ctx, client, jobsCh, resultsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh
		if err != nil {
			return nil, err
		}

		result := <-resultsCh
		res = append(res, result)
	}

	return res, nil
}

func (u *filesUsecase) DeleteFileInGCP(req []*filespkg.DeleteFileReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	for i := range req {
		o := client.Bucket(u.cfg.App().GCPBucket()).Object(req[i].Destination)

		attrs, err := o.Attrs(ctx)
		if err != nil {
			return fmt.Errorf("object.Attrs: %v", err)
		}
		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

		if err := o.Delete(ctx); err != nil {
			return fmt.Errorf("Object(%q).Delete: %v", req[i].Destination, err)
		}
		log.Printf("%v deleted.\n", req[i].Destination)
	}

	return nil
}
