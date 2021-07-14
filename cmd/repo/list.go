package repoCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stewproject/stew/internals/config"
	"github.com/stewproject/stew/util/style"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all repositories.",
	Example: "stew repo list",
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range config.Config.Repositories.Locations {
			fmt.Printf("%s/%s\n", style.Repo.Render(v.Author), style.Repo.Render(v.Name))
		}

		fmt.Printf("\nTotal of %v repositories found.\n", len(config.Config.Repositories.Locations))
	},
}