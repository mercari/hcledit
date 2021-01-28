package command

import (
	"fmt"
	"strconv"

	"github.com/mercari/hcledit"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	Type string
}

func NewCmdCreate() *cobra.Command {
	opts := &CreateOptions{}
	cmd := &cobra.Command{
		Use:   "create <query> <value> <file>",
		Short: "Create a new field",
		Long:  `Runs an address query on a hcl file and create new field with given value.`,
		Args:  cobra.ExactArgs(3),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := runCreate(opts, args); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.Type, "type", "t", "string", "Type of the value")
	return cmd
}

func runCreate(opts *CreateOptions, args []string) error {
	query := args[0]
	valueStr := args[1]
	filePath := args[2]

	editor, err := hcledit.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to read file: %s\n", err)
	}

	value, err := convert(valueStr, opts.Type)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to convert input to specific type: %s\n", err)
	}

	if err := editor.Create(query, value); err != nil {
		return fmt.Errorf("[ERROR] Failed to create: %s\n", err)
	}

	return editor.OverWriteFile()
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
		return nil, fmt.Errorf("Unsupported type: %s", typeStr)
	}
}
