package command

import (
	"github.com/spf13/cobra"
)

func NewCmdRoot(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "hcledit <command> <subcommand> [flags]",
		Short:         "",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.SetVersionTemplate(version)

	cmd.AddCommand(NewCmdVersion(version))
	cmd.AddCommand(NewCmdRead())
	cmd.AddCommand(NewCmdCreate())

	return cmd
}
