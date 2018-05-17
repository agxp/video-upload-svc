package main

import (
	"github.com/minio/minio-go"
	"log"
	"os"
	"time"
	"crypto/md5"
	"encoding/hex"
	"database/sql"
)

type Repository interface {
	S3Request(filename string) (string, error)
	WriteVideoProperties(filename string, title string, description string) (string, string, error)
	UploadFinish(id string) (error)
}

type UploadRepository struct {
	s3 *minio.Client
	pg *sql.DB
}

// should return string
func (repo *UploadRepository) S3Request(filename string) (string, error) {
	log.SetOutput(os.Stdout)
	log.Printf("%#v\n", "filename: " + filename)

	presignedURL, err := repo.s3.PresignedPutObject("videos", filename, time.Hour*24)
	if err != nil {
		log.Printf("%#v\n", "FAILED: with filename: " + filename)
		log.Fatal(err)
		return "", err
	}

	log.Print(presignedURL)

	return presignedURL.String(), nil
}

func (repo *UploadRepository) WriteVideoProperties(filename string, title string, description string) (string, string, error) {
	log.SetOutput(os.Stdout)
	log.Printf("%#v\n", "filename: " + filename)

	objectName := time.Now().String() + filename
	hash := md5.Sum([]byte(objectName))
	id := hex.EncodeToString(hash[:])
	filePath := id + "/" + filename

	insertQuery := `
	INSERT INTO videos(
	id, title, description, date_uploaded, uploaded, 
	date_generated, timeout_date, file_path, view_count, 
	likes, dislikes)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	now := time.Now()

	_, err := repo.pg.Exec(insertQuery, id, title, description, now, false, now,
							now.Add(time.Hour * 24), filePath, 0, 0, 0)

	if err != nil {
		log.Fatal(err)
		return "", "", err
	}


	return id, filePath, nil
}

func (repo *UploadRepository) UploadFinish(id string) error {
	updateQuery := `UPDATE videos SET uploaded=true WHERE id=$1 and timeout_date > $2`

	_, err := repo.pg.Exec(updateQuery, id, time.Now())
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}