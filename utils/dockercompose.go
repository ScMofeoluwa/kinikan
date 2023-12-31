package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DockerCompose struct {
	Services map[string]struct {
		Image string `yaml:"image"`
	} `yaml:"services"`
}

func ExtractImagesFromCompose(filepath ...string) ([]string, error) {
	var path string
	var err error
	if len(filepath) == 0 {
		path, err = findBaseFile()
		if err != nil {
			return nil, err
		}
	} else {
		path = filepath[0]
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read docker-compose file: %w", err)
	}

	var dc DockerCompose
	err = yaml.Unmarshal(data, &dc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal docker-compose file: %w", err)
	}

	var images []string
	for _, service := range dc.Services {
		images = append(images, service.Image)
	}

	return images, nil
}

func findBaseFile() (string, error) {
	filepath := "./docker-compose.yml"
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return "", fmt.Errorf("no docker-compose.yml file found in the root directory")
	}
	return filepath, nil
}
