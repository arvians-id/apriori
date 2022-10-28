package handler

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
)

type StorageValue struct {
	File        *bytes.Reader
	FileName    string
	Size        int64
	ContentType string
}

type StorageService struct {
	AccessKeyID     string
	SecretAccessKey string
	MyRegion        string
	MyBucket        string
}

func NewStorageService() *StorageService {
	return &StorageService{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_KEY"),
		MyRegion:        os.Getenv("AWS_REGION"),
		MyBucket:        os.Getenv("AWS_BUCKET"),
	}
}

func (consumer *StorageService) ConnectToAWS() (*session.Session, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(consumer.MyRegion),
			Credentials: credentials.NewStaticCredentials(consumer.AccessKeyID, consumer.SecretAccessKey, ""),
		},
	)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (consumer *StorageService) UploadToAWS(message *nsq.Message) error {
	sess, err := consumer.ConnectToAWS()
	if err != nil {
		log.Println("[StorageService][UploadToAWS] problem in connect to aws, err: ", err.Error())
		return err
	}

	var storageValue StorageValue
	err = json.Unmarshal(message.Body, &storageValue)
	if err != nil {
		log.Println("[StorageService][Unmarshal] unable to unmarshal data, err: ", err.Error())
		return err
	}
	log.Println("[StorageService][UploadToAWS] storageValue: ", storageValue)

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(consumer.MyBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(storageValue.FileName),
		Body:   storageValue.File,
	})

	if err != nil {
		log.Println("[StorageService][UploadToAWS] problem in upload to aws, err: ", err.Error())
		return err
	}

	return nil
}
