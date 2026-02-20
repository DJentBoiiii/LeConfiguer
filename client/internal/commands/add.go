package commands

import (
	"client/internal/client"

	"github.com/spf13/cobra"
)

func NewAddCommand(apiURL string) *cobra.Command {
	var id, name, configType, environment string
	return &cobra.Command{
		Use:   "add [file-path]",
		Short: "Upload config to MinIO",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireAll(id, name, configType, environment); err != nil {
				return err
			}
			data, err := client.NewClient(apiURL).Upload("/configs", args[0], map[string]string{
				"id": id, "name": name, "type": configType, "environment": environment,
			})
			if err == nil {
				printResult("Config uploaded", data)
			}
			return err
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().StringVarP(&id, "id", "i", "", "Config ID")
			cmd.Flags().StringVarP(&name, "name", "n", "", "Config name")
			cmd.Flags().StringVarP(&configType, "type", "t", "", "Config type")
			cmd.Flags().StringVarP(&environment, "environment", "e", "", "Environment")
		},
	}
}
