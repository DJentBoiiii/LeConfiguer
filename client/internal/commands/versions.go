package commands

import (
	"client/internal/client"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type version struct {
	ID          uint   `json:"id"`
	ConfigID    string `json:"config_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Environment string `json:"environment"`
	Action      string `json:"action"`
	CreatedAt   string `json:"created_at"`
}

func NewVersionsCommand(apiURL string) *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "versions",
		Short: "List all versions of config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireID(id); err != nil {
				return err
			}
			data, err := client.NewClient(apiURL).Get(fmt.Sprintf("/configs/%s/versions", id))
			if err != nil {
				return err
			}
			var versions []version
			if err := json.Unmarshal(data, &versions); err != nil {
				return err
			}
			if len(versions) == 0 {
				fmt.Println("No versions found")
				return nil
			}
			fmt.Printf("%-3s %-12s %-10s %-15s %-20s\n", "ID", "Action", "Type", "Environment", "Created At")
			fmt.Println("---", "----------", "--------", "----------", "------------------")
			for _, v := range versions {
				fmt.Printf("%-3d %-12s %-10s %-15s %-20s\n", v.ID, v.Action, v.Type, v.Environment, v.CreatedAt)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&id, "id", "i", "", "Config ID")
	return cmd
}
