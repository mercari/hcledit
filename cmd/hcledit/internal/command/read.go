package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	"go.mercari.io/hcledit"
)

type ReadOptions struct {
	OutputFormat string
	Fallback     bool
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

	cmd.Flags().StringVarP(&opts.OutputFormat, "output-format", "o", "go-template='{{.Value}}'", "format to print the value as")
	cmd.Flags().BoolVar(&opts.Fallback, "fallback", false, "falls back to reading the raw value if it cannot be evaluated")

	return cmd
}

func runRead(opts *ReadOptions, args []string) (string, error) {
	query, filePath := args[0], args[1]

	editor, err := hcledit.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %s", err)
	}

	readOpts := []hcledit.Option{}
	if opts.Fallback {
		readOpts = append(readOpts, hcledit.WithReadFallbackToRawString())
	}
	results, err := editor.Read(query, readOpts...)
	if err != nil && !opts.Fallback {
		return "", fmt.Errorf("failed to read file: %s", err)
	}

	if strings.HasPrefix(opts.OutputFormat, "go-template") {
		return displayTemplate(opts.OutputFormat, results)
	}

	switch opts.OutputFormat {
	case "json":
		j, err := json.Marshal(results)
		return string(j), err
	case "yaml":
		y, err := yaml.Marshal(results)
		return string(y), err
	default:
		return "", fmt.Errorf("invalid output-format: %s", opts.OutputFormat)
	}
}

func displayTemplate(format string, results map[string]interface{}) (string, error) {
	split := strings.SplitN(format, "=", 2)

	if len(split) != 2 {
		return "", errors.New("go-template should be passed as go-template='<TEMPLATE>'")
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
}
