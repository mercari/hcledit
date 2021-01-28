package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdVersion(version string) *cobra.Command {
	return &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
}
