package commands

import (
	"client/internal/client"
	"fmt"

	"github.com/spf13/cobra"
)

func NewUpdateCommand(apiURL string) *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "update [file-path]",
		Short: "Update config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireID(id); err != nil {
				return err
			}
			c := client.NewClient(apiURL)
			m, err := getMeta(c, id)
			if err != nil {
				return err
			}
			data, err := c.Update(fmt.Sprintf("/configs/%s", id), args[0], map[string]string{
				"id": id, "name": m.Name, "type": m.Type, "environment": m.Environment,
			})
			if err == nil {
				printResult("Config updated", data)
			}
			return err
		},
	}
	cmd.Flags().StringVarP(&id, "id", "i", "", "Config ID")
	return cmd
}
