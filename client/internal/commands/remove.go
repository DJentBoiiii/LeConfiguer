package commands

import (
	"client/internal/client"
	"fmt"

	"github.com/spf13/cobra"
)

func NewRemoveCommand(apiURL string) *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Delete config from MinIO and indexing",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireID(id); err != nil {
				return err
			}
			data, err := client.NewClient(apiURL).Delete(fmt.Sprintf("/configs/%s", id))
			if err == nil {
				printResult("Config deleted", data)
			}
			return err
		},
	}
	cmd.Flags().StringVarP(&id, "id", "i", "", "Config ID")
	return cmd
}
