package service

import (
	"fmt"
	"github.com/arvians-id/apriori/config"
	"github.com/arvians-id/apriori/helper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"strings"
	"sync"
)

type StorageServiceImpl struct {
	AccessKeyID     string
	SecretAccessKey string
	MyRegion        string
	MyBucket        string
}

func NewStorageService(cofiguration config.Config) StorageService {
	return &StorageServiceImpl{
		AccessKeyID:     cofiguration.Get("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: cofiguration.Get("AWS_SECRET_KEY"),
		MyRegion:        cofiguration.Get("AWS_REGION"),
		MyBucket:        cofiguration.Get("AWS_BUCKET"),
	}
}

func (service *StorageServiceImpl) ConnectToAWS() (*session.Session, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(service.MyRegion),
			Credentials: credentials.NewStaticCredentials(service.AccessKeyID, service.SecretAccessKey, ""),
		},
	)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (service *StorageServiceImpl) UploadFile(c *gin.Context, image *multipart.FileHeader) (chan string, error) {
	newFileName := make(chan string)
	go func() {
		extension := strings.Split(image.Filename, ".")
		newFileNames := helper.RandomString(10) + "." + extension[len(extension)-1]

		path, err := helper.GetPath("/assets/", newFileNames)
		if err != nil {
			log.Fatal(err)
		}
		err = c.SaveUploadedFile(image, path)
		if err != nil {
			log.Fatal(err)
		}
		newFileName <- newFileNames
	}()

	return newFileName, nil
}

func (service *StorageServiceImpl) WaitUploadFileS3(file multipart.File, header *multipart.FileHeader, wg *sync.WaitGroup) (string, error) {
	fileName := make(chan string)
	wg.Add(1)
	go func() {
		defer wg.Done()

		sess, err := service.ConnectToAWS()
		if err != nil {
			log.Fatal(err)
		}
		headerFileName := strings.Split(header.Filename, ".")
		fileNames := helper.RandomString(10) + "." + headerFileName[len(headerFileName)-1]
		fileName <- fileNames

		_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
			Bucket:               aws.String(service.MyBucket),
			ACL:                  aws.String("public-read"),
			Key:                  aws.String(fileNames),
			Body:                 file,
			ContentType:          aws.String(header.Header.Get("Content-Type")),
			ContentDisposition:   aws.String("attachment"),
			ServerSideEncryption: aws.String("AES256"),
		})

		if err != nil {
			log.Fatal(err)
		}
	}()

	filePath := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", service.MyBucket, service.MyRegion, <-fileName)
	return filePath, nil
}

func (service *StorageServiceImpl) UploadFileS3(file multipart.File, header *multipart.FileHeader) (string, error) {
	fileName := make(chan string)
	go func() {
		sess, err := service.ConnectToAWS()
		if err != nil {
			log.Fatal(err)
		}
		headerFileName := strings.Split(header.Filename, ".")
		fileNames := helper.RandomString(10) + "." + headerFileName[len(headerFileName)-1]
		fileName <- fileNames

		_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
			Bucket:               aws.String(service.MyBucket),
			ACL:                  aws.String("public-read"),
			Key:                  aws.String(fileNames),
			Body:                 file,
			ContentType:          aws.String(header.Header.Get("Content-Type")),
			ContentDisposition:   aws.String("attachment"),
			ServerSideEncryption: aws.String("AES256"),
		})

		if err != nil {
			log.Fatal(err)
		}
	}()

	filePath := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", service.MyBucket, service.MyRegion, <-fileName)
	return filePath, nil
}

//func (service *StorageServiceImpl) DeleteFileS3(fileName string) error {
//	headerFileName := strings.Split(fileName, "/")
//	oldFileName := headerFileName[len(headerFileName)-1]
//	if oldFileName == "no-image.png" {
//		return nil
//	}
//
//	go func() {
//		sess, err := service.ConnectToAWS()
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		svc := s3.New(sess)
//
//		_, err = svc.DeleteObject(&s3.DeleteObjectInput{
//			Bucket: aws.String(service.MyBucket),
//			Key:    aws.String(oldFileName),
//		})
//		if err != nil {
//			log.Fatal(err)
//		}
//	}()
//
//	return nil
//}
//
//func (service *StorageServiceImpl) WaitDeleteFileS3(oldFileName string, wg *sync.WaitGroup) error {
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		sess, err := service.ConnectToAWS()
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		svc := s3.New(sess)
//
//		_, err = svc.DeleteObject(&s3.DeleteObjectInput{
//			Bucket: aws.String(service.MyBucket),
//			Key:    aws.String(oldFileName),
//		})
//		if err != nil {
//			log.Fatal(err)
//		}
//	}()
//
//	return nil
//}
