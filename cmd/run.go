/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"kinikan/platform"
	"kinikan/platform/heroku"
	"kinikan/platform/railway"
	"kinikan/utils"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	choice   string
	filePath string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run kinikan cli",
	Run: func(cmd *cobra.Command, args []string) {
		var serviceImages []string
		var err error

		s := spinner.New(spinner.CharSets[2], 100*time.Millisecond)

		if filePath != "" {
			s.Suffix = " Fetching Docker compose file..."
			s.Start()

			serviceImages, err = utils.ExtractImagesFromCompose(filePath)
			s.Stop()
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			serviceImages, err = utils.ExtractImagesFromCompose()
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		if choice == "" {
			prompt := promptui.Select{
				Label: "Select PaaS provider",
				Items: []string{"Heroku", "Railway"},
			}

			var err error
			_, choice, err = prompt.Run()
			if err != nil {
				fmt.Printf("prompt failed %v\n", err)
				return
			}
		}

		s.Suffix = " Fetching API Key..."
		s.Start()

		var apiKey string
		var exists bool

		choice = strings.ToLower(choice)
		switch choice {
		case "heroku":
			apiKey, exists = os.LookupEnv("HEROKU_API_KEY")
		case "railway":
			apiKey, exists = os.LookupEnv("RAILWAY_API_KEY")
		}

		if !exists {
			s.FinalMSG = "API key not found\n"
		} else {
			s.FinalMSG = "API key fetched successfully\n"
		}

		s.Stop()
		if !exists {
			return
		}

		platform := initializePlatform(choice, apiKey)

		s.Suffix = " Creating Add-ons..."
		s.Start()
		if err := platform.CreateAddOns(serviceImages); err != nil {
			fmt.Printf("err creating add-ons: %v\n", err)
			s.Stop()
			return
		}
		s.FinalMSG = "Add-ons creation complete"
		s.Stop()
	},
}

func initializePlatform(platform, apiKey string) platform.DeploymentPlatform {
	switch platform {
	case "heroku":
		return heroku.New(apiKey)
	case "railway":
		return railway.New(apiKey)
	default:
		return nil
	}
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&choice, "platform", "p", "", "Specify the PaaS provider")
	runCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "Path to the Docker Compose file")
}
