package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kinikan",
	Short: "Simplify cloud add-on creation with Kinikan CLI.",
	Long: `Kinikan CLI assists developers in automatically setting up necessary add-ons or plugins based on their Docker Compose files. 
	Currently supporting PaaS providers like Heroku and Railway, it simplifies the initial setup process by creating corresponding services for your cloud applications.
	`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
