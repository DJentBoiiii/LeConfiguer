package commands

import (
	"bytes"
	"client/internal/client"
	"fmt"
	"os"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

func NewDiffCommand(apiURL string) *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "diff [file-path]",
		Short: "Diff with previous version if exists",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireID(id); err != nil {
				return err
			}
			localContent, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			remoteContent, err := client.NewClient(apiURL).Get(fmt.Sprintf("/diff/%s", id))
			if err != nil {
				return err
			}
			if bytes.Equal(localContent, remoteContent) {
				fmt.Println("No differences")
				return nil
			}
			text, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
				A:        difflib.SplitLines(string(remoteContent)),
				B:        difflib.SplitLines(string(localContent)),
				FromFile: "remote", ToFile: "local", Context: 3,
			})
			fmt.Print(text)
			return nil
		},
	}
	cmd.Flags().StringVarP(&id, "id", "i", "", "Config ID")
	return cmd
}
