package model

import (
	"bytes"
)

type EmailService struct {
	ToEmail string
	Subject string
	Message string
}

type StorageService struct {
	File        *bytes.Reader
	FileName    string
	Size        int64
	ContentType string
}
