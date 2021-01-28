package command

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/mercari/hcledit"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type ReadOptions struct {
	OutputFormat string
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

	cmd.Flags().StringVarP(&opts.OutputFormat, "output-format", "o", "go-template='{{.Key}} {{.Value}}'", "format to print the value as")

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

	if strings.HasPrefix(opts.OutputFormat, "go-template") {
		split := strings.SplitN(opts.OutputFormat, "=", 2)
		if len(split) != 2 {
			return "", fmt.Errorf(`[ERROR] go-template should be passed as go-template='<TEMPLATE>'`)
		}

		templateFormat := strings.Trim(split[1], "'")

		tmpl, err := template.New("output").Parse(templateFormat)
		if err != nil {
			return "", err
		}

		var result strings.Builder

		for key, value := range results {
			formatted := struct {
				Key   string
				Value string
			}{
				fmt.Sprintf("%v", key),
				fmt.Sprintf("%v", value),
			}

			if err := tmpl.Execute(&result, formatted); err != nil {
				return result.String(), err
			}
		}

		return result.String(), nil
	} else if opts.OutputFormat == "json" {
		j, err := json.Marshal(results)
		return string(j), err
	} else if opts.OutputFormat == "yaml" {
		y, err := yaml.Marshal(results)
		return string(y), err
	}

	return "", fmt.Errorf("[ERROR] Invalid output-format")
}
