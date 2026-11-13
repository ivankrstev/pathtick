package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/ivankrstev/pathtick/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewRootCmd() *cobra.Command {
	options := app.DefaultOptions()

	format := string(options.Format)

	cmd := &cobra.Command{
		Use:           "pathtick [target-directory]",
		Short:         "Export a directory tree as a checklist",
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				options.TargetDir = args[0]
			}

			options.Format = app.OutputFormat(format)

			// Reject explicit file-output options when stdout mode is enabled.
			if options.Stdout {
				if cmd.Flags().Changed("output-dir") {
					return errors.New("--stdout and --output-dir cannot be used together")
				}
				if cmd.Flags().Changed("output-name") {
					return errors.New("--stdout and --output-name cannot be used together")
				}
			}

			return options.Validate()
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(options, cmd.OutOrStdout())
		},
	}

	flags := cmd.Flags()

	// Formatting
	flags.StringVarP(&format, "format", "f", string(options.Format), "Output format (currently only md is supported)")
	flags.StringVarP(&options.Title, "title", "t", options.Title, "Title shown in the exported document (default: derived from the target directory)")

	// Output
	flags.StringVarP(&options.OutputDir, "output-dir", "o", options.OutputDir, "Directory to write the output files")
	flags.StringVarP(&options.OutputBaseName, "output-name", "n", options.OutputBaseName, "Base name for the output files (without extension)")
	flags.BoolVar(&options.Stdout, "stdout", options.Stdout, "Write output to stdout instead of files")
	flags.BoolVar(&options.Overwrite, "overwrite", options.Overwrite, "Overwrite existing output files")

	cmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		out := cmd.OutOrStdout()
		helpFlags := cmd.Flags()

		fmt.Fprintf(out, "%s\n\nUsage:\n  %s\n\n", cmd.Short, cmd.UseLine())

		fmt.Fprintln(out, "Arguments:")
		fmt.Fprintln(out, "  target-directory  Directory to scan (default: current directory)")
		fmt.Fprintln(out)

		printFlagGroup(
			out,
			"Formatting",
			helpFlags,
			"format",
			"title",
		)

		printFlagGroup(
			out,
			"Output",
			helpFlags,
			"output-dir",
			"output-name",
			"stdout",
			"overwrite",
		)

		printFlagGroup(
			out,
			"Help",
			helpFlags,
			"help",
		)
	})

	return cmd
}

// printFlagGroup writes a help section containing the named flags.
// Flags are printed in the order listed, and unknown names are skipped.
func printFlagGroup(
	out io.Writer,
	title string,
	flags *pflag.FlagSet,
	names ...string,
) {
	fmt.Fprintf(out, "%s:\n", title)

	for _, name := range names {
		flag := flags.Lookup(name)
		if flag == nil {
			continue
		}

		shorthand := ""
		if flag.Shorthand != "" {
			shorthand = "-" + flag.Shorthand + ", "
		}

		valueType := ""
		if flag.Value.Type() != "bool" {
			valueType = " " + flag.Value.Type()
		}

		defaultValue := ""
		if flag.DefValue != "" &&
			(flag.Value.Type() != "bool" || flag.DefValue != "false") {
			defaultValue = " (default " + flag.DefValue + ")"
		}

		fmt.Fprintf(out, "  %s--%s%s\t%s%s\n",
			shorthand, flag.Name, valueType, flag.Usage, defaultValue)
	}

	fmt.Fprintln(out)
}
