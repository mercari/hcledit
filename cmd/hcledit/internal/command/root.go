package command

import (
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "hcledit <command> <subcommand> [flags]",
		Short:         "",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		NewCmdVersion(),
		NewCmdRead(),
		NewCmdCreate(),
		NewCmdUpdate(),
		NewCmdDelete(),
	)

	return cmd
}
