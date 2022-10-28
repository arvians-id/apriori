package util

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func OpenCsvFile(path string) ([][]string, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.LazyQuotes = true

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

func GetPathAWS(fileName string) (string, string) {
	headerFileName := strings.Split(fileName, ".")
	randomName := RandomString(10) + "." + headerFileName[len(headerFileName)-1]
	filePath := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), randomName)

	return filePath, randomName
}
