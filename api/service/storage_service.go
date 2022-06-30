package service

import (
	"apriori/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type StorageService interface {
	ConnectToAWS() (*session.Session, error)
}

type storageService struct {
	AccessKeyID     string
	SecretAccessKey string
	MyRegion        string
}

func NewStorageService() StorageService {
	return &storageService{
		AccessKeyID:     config.New().Get("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: config.New().Get("AWS_SECRET_KEY"),
		MyRegion:        config.New().Get("AWS_REGION"),
	}
}

func (service *storageService) ConnectToAWS() (*session.Session, error) {
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
