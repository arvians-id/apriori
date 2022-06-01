package setup

import "github.com/joho/godotenv"

func TestEnv() error {
	err := godotenv.Load("../../.env")
	if err != nil {
		return err
	}

	return nil
}
