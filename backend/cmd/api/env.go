package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func loadLocalEnvironment() error {
	if path := os.Getenv("AEOLYZER_ENV_FILE"); path != "" {
		if err := godotenv.Load(path); err != nil {
			return fmt.Errorf("load environment file: %w", err)
		}
		return nil
	}
	for _, path := range []string{".env", "../.env", "../../.env", "../../../.env"} {
		info, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("inspect environment file: %w", err)
		}
		if info.IsDir() {
			return errors.New("environment file path is a directory")
		}
		if err := godotenv.Load(path); err != nil {
			return fmt.Errorf("load environment file: %w", err)
		}
		return nil
	}
	return nil
}
