package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mercari/hcledit"
)

func NewCmdDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <query> <file>",
		Short: "Delete the given field and its value",
		Long:  `Runs an address query on a hcl file and delete the given field and its value.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := runDelete(args); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func runDelete(args []string) error {
	query := args[0]
	filePath := args[1]

	editor, err := hcledit.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	if err := editor.Delete(query); err != nil {
		return fmt.Errorf("failed to delete: %s", err)
	}

	return editor.OverWriteFile()
}
