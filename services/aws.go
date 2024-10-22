package services

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	BucketName = aws.String(LoadEnv("AWS_BUCKET_NAME"))
	Prefix     = LoadEnv("ENV")
)

func Aws() *s3.Client {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(LoadEnv("AWS_REGION")),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     LoadEnv("AWS_ACCESS_KEY"),
				SecretAccessKey: LoadEnv("AWS_SECRET_KEY"),
			},
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	return client
}

func Upload(f *multipart.FileHeader, id string) error {
	client := Aws()
	ctx := context.TODO()

	file, err := f.Open()
	if err != nil {
		return err
	}

	upload := &s3.PutObjectInput{
		Bucket: BucketName,
		Key:    aws.String(Prefix + "/" + id + "/" + f.Filename),
		Body:   file,
		ACL:    "public-read", // temporário
	}

	_, err = client.PutObject(ctx, upload)

	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("File uploaded successfully")

	return nil
}
