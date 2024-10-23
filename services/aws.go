// Vai conter as funções de upload de arquivos para o S3 da AWS
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

// Variáveis globais para o nome do bucket e o prefixo
var (
	BucketName = aws.String(LoadEnv("AWS_BUCKET_NAME"))
	Prefix     = LoadEnv("ENV")
)

// Função para criar o client da AWS
func Aws() *s3.Client {
	ctx := context.TODO() // context é uma interface vazia, por isso o uso do TODO

	// Configuração do client da AWS, com as credenciais e a região que vamos usar
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

	return client // Retorna o client da AWS
}

func UploadToAws(f *multipart.FileHeader, id string) error {
	client := Aws()
	ctx := context.TODO()

	// Abre o arquivo que foi enviado, já ue o arquivo é um multipart.FileHeader não vai para o S3
	file, err := f.Open()
	if err != nil {
		return err
	}

	// Configuração do upload do arquivo para o S3, definindo o bucket, a key (nome do arquivo) e o corpo do arquivo
	upload := &s3.PutObjectInput{
		Bucket: BucketName,
		Key:    aws.String(Prefix + "/" + id + "/" + f.Filename),
		Body:   file,
		ACL:    "public-read", // temporário
	}

	_, err = client.PutObject(ctx, upload) // Fazendo o upload em si

	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("File uploaded successfully")

	return nil // Retorna nil se tudo ocorreu bem
}
