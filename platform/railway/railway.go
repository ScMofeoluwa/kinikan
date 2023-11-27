package railway

import (
	"fmt"
	"kinikan/platform"
	"kinikan/utils"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

const url = "https://backboard.railway.app/graphql/v2"

var addOns = map[string]string{
	"postgres": "postgresql",
	"mysql":    "mysql",
	"redis":    "redis",
	"mongo":    "mongodb",
}

type CreateProjectResponse struct {
	Data struct {
		ProjectCreate struct {
			ID string `json:"id"`
		} `json:"projectCreate"`
	} `json:"data"`
}

type Railway struct {
	req *resty.Request
}

var _ platform.DeploymentPlatform = &Railway{}

func New(apiKey string) *Railway {
	client := resty.New().
		SetRetryCount(3).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusInternalServerError
			},
		)

	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(apiKey)

	return &Railway{
		req: req,
	}
}

func (r *Railway) CreateApp() (string, error) {
	appName := utils.RandomAppName()

	query := `
  	mutation projectCreate($name: String!) {
  		projectCreate(input: {name: $name}) {
  			id
  		}
  	}
	`

	variables := map[string]string{
		"name": appName,
	}

	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	var resp CreateProjectResponse
	_, err := r.req.SetBody(&payload).SetResult(&resp).Post(url)
	if err != nil {
		return "", fmt.Errorf("error creating Railway project: %v", err)
	}

	projectID := resp.Data.ProjectCreate.ID
	fmt.Printf("\nSuccessfully created Railway project: %s\n", appName)
	return projectID, nil
}

func (r *Railway) CreateAddOns(serviceImages []string) error {
	projectID, err := r.CreateApp()
	if err != nil {
		return err
	}

	query := `
  	mutation pluginCreate($name: String!, $projectId: String!) {
  		pluginCreate(input: { name: $name, projectId: $projectId }) {
  			status
  		}
		}
	`

	for _, serviceImage := range serviceImages {
		baseImg := strings.Split(serviceImage, ":")[0]
		addOn, exists := addOns[baseImg]
		if !exists {
			continue
		}

		variables := map[string]string{
			"name":      addOn,
			"projectId": projectID,
		}

		payload := map[string]interface{}{
			"query":     query,
			"variables": variables,
		}

		_, err := r.req.SetBody(&payload).Post(url)
		if err != nil {
			return fmt.Errorf("error attaching plugin to project: %v", err)
		}
	}

	return nil
}

func (r *Railway) Deploy() error {
	return nil
}
