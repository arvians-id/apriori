package utils

import (
	"apriori/config"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func OpenCsvFile(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	all, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return all, nil
}

func GetPath(path string, file string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filename := path + file
	fullPath := filepath.Join(dir, filename)

	return fullPath, nil
}

func UploadImage(c *gin.Context, image *multipart.FileHeader) (string, error) {
	extension := strings.Split(image.Filename, ".")
	newFileName := RandomString(10) + "." + extension[len(extension)-1]

	path, err := GetPath("/assets/", newFileName)
	if err != nil {
		return "", err
	}
	err = c.SaveUploadedFile(image, path)
	if err != nil {
		return "", err
	}

	return newFileName, nil
}

func UploadToS3(c *gin.Context) (string, error) {
	sess := c.MustGet("sess").(*session.Session)

	configuration := config.New()
	myBucket := configuration.Get("AWS_BUCKET")
	myRegion := configuration.Get("AWS_REGION")

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		return "", err
	}
	headerFileName := strings.Split(header.Filename, ".")
	fileName := RandomString(10) + "." + headerFileName[len(headerFileName)-1]

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(myBucket),
		ACL:                  aws.String("public-read"),
		Key:                  aws.String(fileName),
		Body:                 file,
		ContentType:          aws.String(header.Header.Get("Content-Type")),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", myBucket, myRegion, fileName)
	return filePath, nil
}
