package twigga

import (
	"encoding/json"
	"io"
	"os"
)

func LoadConfig(filename string) (*BongoCloudClient, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config BongoCloudClient
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
