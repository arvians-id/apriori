package aws

import (
	"bytes"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/arvians-id/apriori/cmd/config"
	"github.com/arvians-id/apriori/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

type StorageS3 struct {
	AccessKeyID     string
	SecretAccessKey string
	MyRegion        string
	MyBucket        string
}

func NewStorageS3(cofiguration config.Config) *StorageS3 {
	return &StorageS3{
		AccessKeyID:     cofiguration.Get("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: cofiguration.Get("AWS_SECRET_KEY"),
		MyRegion:        cofiguration.Get("AWS_REGION"),
		MyBucket:        cofiguration.Get("AWS_BUCKET"),
	}
}

func (storageS3 *StorageS3) ConnectToAWS() (*session.Session, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(storageS3.MyRegion),
			Credentials: credentials.NewStaticCredentials(storageS3.AccessKeyID, storageS3.SecretAccessKey, ""),
		},
	)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (storageS3 *StorageS3) UploadToAWS(file multipart.File, fileName string, contentType string) error {
	sess, err := storageS3.ConnectToAWS()
	if err != nil {
		log.Println("[AWS][UploadToAWS][ConnectToAWS] problem in connecting to aws, err: ", err.Error())
		return err
	}

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(storageS3.MyBucket),
		ACL:                  aws.String("public-read"),
		Key:                  aws.String(fileName),
		Body:                 file,
		ContentType:          aws.String(contentType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		log.Println("[AWS][UploadToAWS][PutObject] problem in upload to aws, err: ", err.Error())
		return err
	}

	return nil
}

func (storageS3 *StorageS3) UploadToAWS2(fileUpload graphql.Upload, initFileName string) (string, error) {
	sess, err := storageS3.ConnectToAWS()
	if err != nil {
		log.Println("[AWS][UploadToAWS2][ConnectToAWS] problem in connecting to aws, err: ", err.Error())
		return "", err
	}

	// Read the file
	stream, err := ioutil.ReadAll(fileUpload.File)
	if err != nil {
		log.Println("[AWS][UploadToAWS2][ReadAll] problem in reading file, err: ", err.Error())
		return "", err
	}

	// then write it to a file
	err = ioutil.WriteFile(initFileName, stream, 0644)
	if err != nil {
		log.Println("[AWS][UploadToAWS2][WriteFile] problem in writing file, err: ", err.Error())
		return "", err
	}

	// Open the file
	file, err := os.Open(initFileName)
	if err != nil {
		log.Println("[AWS][UploadToAWS2][Open] problem in opening file, err: ", err.Error())
		return "", err
	}
	defer file.Close()

	// Upload the file to S3.
	buffer := make([]byte, fileUpload.Size)
	_, err = file.Read(buffer)
	if err != nil {
		log.Println("[AWS][UploadToAWS2][Read] problem in reading file, err: ", err.Error())
		return "", err
	}
	fileBytes := bytes.NewReader(buffer)

	headerFileName := strings.Split(fileUpload.Filename, ".")
	fileNames := util.RandomString(10) + "." + headerFileName[len(headerFileName)-1]

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(storageS3.MyBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(fileNames),
		Body:   fileBytes,
	})
	if err != nil {
		log.Println("[AWS][UploadToAWS2][PutObject] problem in upload to aws, err: ", err.Error())
		return "", err
	}

	filePath := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", storageS3.MyBucket, storageS3.MyRegion, fileNames)
	return filePath, nil
}

// Deprecated on production
func (storageS3 *StorageS3) DeleteFromAWS(filePath string) error {
	headerFilePathName := strings.Split(filePath, "/")
	fileName := headerFilePathName[len(headerFilePathName)-1]
	if fileName == "no-image.png" {
		return nil
	}

	sess, err := storageS3.ConnectToAWS()
	if err != nil {
		log.Println("[AWS][DeleteFromAWS][ConnectToAWS] problem in connecting to aws, err: ", err.Error())
		return err
	}

	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(storageS3.MyBucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		log.Println("[AWS][DeleteFromAWS][DeleteObject] problem in delete from aws, err: ", err.Error())
		return err
	}

	return nil
}
