package utils

import (
	"encoding/csv"
	"os"
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
