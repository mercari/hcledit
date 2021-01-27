package command

import (
	"fmt"
	"strings"

	"github.com/mercari/hcledit"
	"github.com/spf13/cobra"
)

type ReadOptions struct {
	ValueFormat string
	ValueOnly   bool
}

func NewCmdRead() *cobra.Command {
	opts := &ReadOptions{}
	cmd := &cobra.Command{
		Use:   "read <query> <file>",
		Short: "Read a value",
		Long:  `Runs an address query on a hcl file and prints the result`,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			result, err := runRead(opts, args)
			if err != nil {
				return err
			}

			fmt.Print(result)
			return nil
		},
	}

	cmd.Flags().BoolVar(&opts.ValueOnly, "value-only", false, "only print the value")
	cmd.Flags().StringVar(&opts.ValueFormat, "value-format", "%v", "format to print the value as")

	return cmd
}

func runRead(opts *ReadOptions, args []string) (string, error) {
	query := args[0]
	filePath := args[1]

	editor, err := hcledit.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Failed to read file: %s\n", err)
	}

	results, err := editor.Read(query)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Failed to read file: %s\n", err)
	}

	var result strings.Builder

	for key, value := range results {
		if opts.ValueOnly {
			fmt.Fprintf(&result, opts.ValueFormat, value)
		} else {
			fmt.Fprintln(&result, key, value)
		}
	}

	return result.String(), nil
}
