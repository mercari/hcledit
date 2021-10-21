package command

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"go.mercari.io/hcledit"
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

func convert(inputStr, typeStr string) (interface{}, error) {
	switch typeStr {
	case "string":
		return inputStr, nil
	case "int":
		return strconv.Atoi(inputStr)
	case "bool":
		return strconv.ParseBool(inputStr)
	case "raw":
		return hcledit.RawVal(inputStr), nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", typeStr)
	}
}
