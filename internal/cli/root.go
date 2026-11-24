package cli

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/ivankrstev/pathtick/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewRootCmd() *cobra.Command {
	options := app.DefaultOptions()

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

			// Reject explicit file-output options when stdout mode is enabled.
			if options.Output.Stdout {
				for _, name := range []string{"output-dir", "output-name", "overwrite"} {
					if cmd.Flags().Changed(name) {
						return fmt.Errorf("--stdout and --%s cannot be used together", name)
					}
				}
			}

			return options.Validate()
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(options, cmd.OutOrStdout())
		},
	}

	flags := cmd.Flags()

	// Enum flags use VarP: the enum types implement pflag.Value, so they bind
	// directly to the option fields. Bad values are rejected while parsing,
	// and the help text shows the allowed values via Type().

	// Scan
	flags.IntVarP(&options.Scan.Depth, "depth", "d", options.Scan.Depth, "Maximum directory traversal depth; -1 means unlimited")
	flags.BoolVar(&options.Scan.SkipSymlinks, "skip-symlinks", options.Scan.SkipSymlinks, "Exclude symbolic link entries")
	flags.BoolVar(&options.Scan.SkipHidden, "skip-hidden", options.Scan.SkipHidden, "Exclude entries whose names start with a dot")

	// Display
	flags.StringVarP(&options.Display.Title, "title", "t", options.Display.Title, "Title shown in the exported document (default: derived from the target directory)")
	flags.VarP(&options.Display.Sort, "sort", "s", "Sort directory entries, case-insensitively: none, asc, or desc")

	// Output
	flags.VarP(&options.Output.Format, "format", "f", "Output format")
	flags.StringVarP(&options.Output.Dir, "output-dir", "o", options.Output.Dir, "Directory to write the output files")
	flags.StringVarP(&options.Output.BaseName, "output-name", "n", options.Output.BaseName, "Base name for the output files (without extension)")
	flags.BoolVar(&options.Output.Stdout, "stdout", options.Output.Stdout, "Write output to stdout instead of files")
	flags.BoolVar(&options.Output.Overwrite, "overwrite", options.Output.Overwrite, "Overwrite existing output files")

	cmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		out := cmd.OutOrStdout()
		helpFlags := cmd.Flags()

		fmt.Fprintf(out, "%s\n\nUsage:\n  %s\n\n", cmd.Short, cmd.UseLine())

		fmt.Fprintln(out, "Arguments:")
		fmt.Fprintln(out, "  target-directory  Directory to scan (default: current directory)")
		fmt.Fprintln(out)

		printFlagGroup(
			out,
			"Scan",
			helpFlags,
			"depth",
			"skip-symlinks",
			"skip-hidden",
		)

		printFlagGroup(
			out,
			"Display",
			helpFlags,
			"title",
			"sort",
		)

		printFlagGroup(
			out,
			"Output",
			helpFlags,
			"format",
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
	w := tabwriter.NewWriter(out, 0, 0, 4, ' ', 0)

	for _, name := range names {
		flag := flags.Lookup(name)
		if flag == nil {
			continue
		}

		shorthand := "    "
		if flag.Shorthand != "" {
			shorthand = "-" + flag.Shorthand + ", "
		}

		// The value name is the `backquoted` word in the usage text if there
		// is one, otherwise the flag's type ("int", "md|html|both", ...).
		valueName, usage := pflag.UnquoteUsage(flag)
		if valueName != "" {
			valueName = " " + valueName
		}

		defaultValue := ""
		if flag.DefValue != "" && flag.DefValue != "false" && flag.DefValue != "[]" {
			defaultValue = " (default " + flag.DefValue + ")"
		}

		fmt.Fprintf(w, "  %s--%s%s\t%s%s\n",
			shorthand, flag.Name, valueName, usage, defaultValue)
	}

	_ = w.Flush()
	fmt.Fprintln(out)
}
