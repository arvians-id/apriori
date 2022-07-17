package utils

import (
	"encoding/csv"
	"net/http"
	"os"
	"path/filepath"
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
