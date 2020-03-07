package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "helm push-all <directory-with-charts> <chartmuseum-repo>",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Helm plugin to push all charts in a directory to chartmuseum",
	Long: `Helm plugin to push all charts in a directory to chartmuseum

Examples:
  $ helm push-all all-my-charts chartmuseum                 # push using chart repo name
  $ helm push-all all-my-charts https://my.chart.repo.com   # push directly to chart repo URL
  $ helm push-all . https://my.chart.repo.com               # push current diractory containing all charts
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
