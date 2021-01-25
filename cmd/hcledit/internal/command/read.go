package command

import (
	"fmt"

	"github.com/mercari/hcledit"
	"github.com/posener/complete"
)

type readCommand struct{}

func (c *readCommand) Synopsis() string {
	return "merctl read <QUERY> <FILE>"
}

func (c *readCommand) Help() string {
	return "read"
}

func (c *readCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}
func (c *readCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *readCommand) Run(args []string) int {
	if len(args) != 2 {
		fmt.Printf("[ERROR] Invalid arguments\n")
		return 1
	}

	query := args[0]
	filePath := args[1]

	editor, err := hcledit.ReadFile(filePath)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read file: %s\n", err)
		return 1
	}

	results, err := editor.Read(query)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read file: %s\n", err)
		return 1
	}

	for key, value := range results {
		fmt.Println(key, value)
	}

	return 0
}
