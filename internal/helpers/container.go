package helpers

import (
	"io"
	"os"
)

var (
	ContainerId string
)

func StoreContainerId(id string) error {
	if _, err := GetContainerId(); err != nil {
		file, err := os.OpenFile("container_id", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}

		file.WriteString(id)
		file.Close()
	}
	return nil
}

func GetContainerId() (string, error) {
	if ContainerId == "" {
		file, err := os.OpenFile("container_id", os.O_RDWR, 0644)
		if err != nil {
			return "", err
		}
		allBytes, err := io.ReadAll(file)
		if err != nil {
			return "", err
		}

		ContainerId = string(allBytes)
		defer file.Close()
	}

	return ContainerId, nil
}
