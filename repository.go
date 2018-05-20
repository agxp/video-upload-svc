package main

import (
	"github.com/minio/minio-go"
	"log"
	"os"
	"time"
	"crypto/md5"
	"encoding/hex"
	"database/sql"
	"github.com/opentracing/opentracing-go"
	"context"
	"go.uber.org/zap"
)

type Repository interface {
	S3Request(p opentracing.SpanContext, filename string) (string, error)
	WriteVideoProperties(p opentracing.SpanContext, filename string, title string, description string) (string, string, error)
	UploadFinish(p opentracing.SpanContext, id string) (error)
}

type UploadRepository struct {
	s3 *minio.Client
	pg *sql.DB
	tracer *opentracing.Tracer
}

// should return string
func (repo *UploadRepository) S3Request(parent opentracing.SpanContext, filename string) (string, error) {
	sp, _ := opentracing.StartSpanFromContext(context.Background(), "S3Request_Repo", opentracing.ChildOf(parent))

	sp.LogKV("filename", filename)


	defer sp.Finish()
	log.SetOutput(os.Stdout)
	log.Printf("%#v\n", "filename: " + filename)

	psSP, _ := opentracing.StartSpanFromContext(context.Background(), "S3_PresignedPutObject", opentracing.ChildOf(sp.Context()))

	psSP.LogKV("filename", filename)

	presignedURL, err := repo.s3.PresignedPutObject("videos", filename, time.Hour*24)
	if err != nil {
		log.Printf("%#v\n", "FAILED: with filename: " + filename)
		log.Fatal(err)
		psSP.Finish()
		return "", err
	}
	psSP.Finish()

	log.Print(presignedURL)

	return presignedURL.String(), nil
}

func (repo *UploadRepository) WriteVideoProperties(p opentracing.SpanContext, filename string, title string, description string) (string, string, error) {
	sp, _ := opentracing.StartSpanFromContext(context.Background(), "WriteVideoProperties_Repo", opentracing.ChildOf(p))

	sp.LogKV("filename", filename,"title", title,"description", description)

	defer sp.Finish()

	logger.Info("filename", zap.String("filename", filename))

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

	dbSP, _ := opentracing.StartSpanFromContext(context.Background(),"PG_WriteVideoProperties", opentracing.ChildOf(sp.Context()))

	dbSP.LogKV("id", id, "title", title, "description", description, "filePath", filePath)

	_, err := repo.pg.Exec(insertQuery, id, title, description, now, false, now,
							now.Add(time.Hour * 24), filePath, 0, 0, 0)

	if err != nil {
		log.Fatal(err)
		dbSP.Finish()
		return "", "", err
	}
	dbSP.Finish()


	return id, filePath, nil
}

func (repo *UploadRepository) UploadFinish(p opentracing.SpanContext, id string) error {
	sp, _ := opentracing.StartSpanFromContext(context.Background(), "UploadFinish_Repo", opentracing.ChildOf(p))

	sp.LogKV("id", id)

	defer sp.Finish()
	updateQuery := `UPDATE videos SET uploaded=true WHERE id=$1 and timeout_date > $2`

	dbSP, _ := opentracing.StartSpanFromContext(context.Background(), "PG_UploadFinish", opentracing.ChildOf(sp.Context()))
	dbSP.LogKV("updateQuery", updateQuery, "id", id)

	_, err := repo.pg.Exec(updateQuery, id, time.Now())
	if err != nil {
		log.Fatal(err)
		dbSP.Finish()
		return err
	}
	dbSP.Finish()
	return nil
}