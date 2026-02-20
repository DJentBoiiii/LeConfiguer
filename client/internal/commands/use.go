package commands

import (
	"client/internal/client"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewUseCommand(apiURL string) *cobra.Command {
	var id, outputFile string
	return &cobra.Command{
		Use:   "use",
		Short: "Download file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireID(id); err != nil {
				return err
			}
			data, err := client.NewClient(apiURL).Get(fmt.Sprintf("/configs/%s", id))
			if err != nil {
				return err
			}
			if outputFile != "" {
				os.WriteFile(outputFile, data, 0644)
				fmt.Fprintf(os.Stderr, "Downloaded to: %s\n", outputFile)
			} else {
				fmt.Print(string(data))
			}
			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().StringVarP(&id, "id", "i", "", "Config ID")
			cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file")
		},
	}
}
