package heroku

import (
	"fmt"
	"os/exec"
	"strings"
)

var addOns = map[string]string{
	"postgres": "heroku-postgresql:hobby-dev",
	"mysql":    "jawsdb:kitefin",
	"redis":    "rediscloud:30",
}

type Heroku struct{}

func New() *Heroku {
	return &Heroku{}
}

func (h *Heroku) ValidateCLI() error {
	_, err := exec.LookPath("heroku")
	if err != nil {
		return fmt.Errorf("Heroku CLI is not installed")
	}

	cmd := exec.Command("heroku", "whoami")
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("not logged in to Heroku. Please log in using 'heroku login'. Error: %s", err)
	}

	return nil
}

func (h *Heroku) CreateApp() (string, error) {
	cmd := exec.Command("heroku", "apps:create")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error creating Heroku app: %v - Output: %s", err, string(out))
	}

	appName := strings.Split(strings.Split(string(out), "//")[1], ".")[0]
	parts := strings.Split(appName, "-")
	appName = strings.Join(parts[:len(parts)-1], "-")
	fmt.Printf("Successfully created heroku app: %s\n", appName)
	return appName, nil
}

func (h *Heroku) CreateAddOns(serviceImages []string) error {
	appName, err := h.CreateApp()
	if err != nil {
		return err
	}

	cmd := exec.Command("heroku", "git:remote", "-a", appName)
	var out []byte
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error linking app to heroku: %v - Output: %s", err, string(out))
	}
	fmt.Printf("Successfully linked app: %s to Heroku \n", appName)

	for _, serviceImage := range serviceImages {
		baseImg := strings.Split(serviceImage, ":")[0]
		addOn, exists := addOns[baseImg]
		if !exists {
			fmt.Printf("No Heroku add-on for the Docker image: %s\n", serviceImage)
			continue
		}

		cmd := exec.Command("heroku", "addons:create", addOn)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error creating Heroku add-on for %s: %v - Output: %s", serviceImage, err, string(out))
		}

		fmt.Printf("Successfully created Heroku add-on %s\n", addOn)
	}

	return nil
}
