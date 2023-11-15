package heroku

import (
	"fmt"
	"kinikan/platform"
	"net/http"
	"os/exec"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	appUrl   = "https://api.heroku.com/apps"
	addOnUrl = "https://api.heroku.com/apps/%s/addons"
)

var addOns = map[string]string{
	"postgres": "heroku-postgresql:hobby-dev",
	"mysql":    "jawsdb:kitefin",
	"redis":    "rediscloud:30",
}

type CreateAppResponse struct {
	Name   string `json:"name"`
	GitURL string `json:"git_url"`
}

type Heroku struct {
	req *resty.Request
}

var _ platform.DeploymentPlatform = &Heroku{}

func New(apiKey string) *Heroku {
	client := resty.New().
		SetRetryCount(3).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusInternalServerError
			},
		)

	req := client.R().
		SetHeaders(
			map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/vnd.heroku+json; version=3",
			},
		).SetAuthToken(apiKey)

	return &Heroku{
		req: req,
	}
}

func (h *Heroku) CreateApp() (string, error) {
	var resp CreateAppResponse
	_, err := h.req.SetResult(&resp).Post(appUrl)
	if err != nil {
		return "", fmt.Errorf("error creating Heroku app: %v", err)
	}

	out, err := exec.Command("git", "remote", "-v").Output()
	if err != nil {
		return "", err
	}

	if strings.Contains(string(out), "heroku") {
		exec.Command("git", "remote", "set-url", "heroku", resp.GitURL).Run()
	} else {
		exec.Command("git", "remote", "add", "heroku", resp.GitURL).Run()
	}

	fmt.Printf("Successfully created Railway project: %s\n", resp.Name)
	return resp.Name, nil
}

func (h *Heroku) CreateAddOns(serviceImages []string) error {
	appName, err := h.CreateApp()
	if err != nil {
		return err
	}

	for _, serviceImage := range serviceImages {
		baseImg := strings.Split(serviceImage, ":")[0]
		addOn, exists := addOns[baseImg]
		if !exists {
			fmt.Printf("No Heroku add-on for the Docker image: %s\n", serviceImage)
			continue
		}

		payload := map[string]string{
			"plan": addOn,
		}

		_, err := h.req.SetBody(payload).Post(fmt.Sprintf(addOnUrl, appName))
		if err != nil {
			return fmt.Errorf("error attaching add-on to app: %v", err)
		}

		fmt.Printf("Successfully created Heroku add-on %s\n", addOn)
	}

	return nil
}

func (h *Heroku) Deploy() error {
	return nil
}
