package commands

import (
	"client/internal/client"
	"fmt"

	"github.com/spf13/cobra"
)

func NewRollbackCommand(apiURL string) *cobra.Command {
	var id string
	return &cobra.Command{
		Use:   "rollback",
		Short: "Rollback config to previous version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireID(id); err != nil {
				return err
			}
			data, err := client.NewClient(apiURL).Post(fmt.Sprintf("/configs/%s/rollback", id))
			if err == nil {
				printResult("Rolled back", data)
			}
			return err
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().StringVarP(&id, "id", "i", "", "Config ID")
		},
	}
}
