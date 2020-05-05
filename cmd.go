package orlop

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/switch-bit/orlop/log"
)

// RootCommand represents the root command object
type RootCommand struct {
	c *cobra.Command
}

// Execute the command
func (r RootCommand) Execute() {
	r.c.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "help about any command",
		Long: `Help provides help for any command in the application.
Simply type ` + r.c.Name() + ` help [path to command] for full details.`,

		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				c.Printf("Unknown help topic %#q\n", args)
				c.Root().Usage()
			} else {
				cmd.InitDefaultHelpFlag() // make possible 'help' flag to be shown
				cmd.Help()
			}
		},
	})

	r.c.PersistentFlags().String("config", fmt.Sprintf(".%s/config.yaml", r.c.Name()), "configuration file to use")
	r.c.PersistentFlags().String("loglevel", "", "logging level")

	if err := r.c.Execute(); err != nil {
		log.Fatal(err)
	}
}

// AddCommand adds a sub command
func (r RootCommand) AddCommand(c *cobra.Command) {
	r.c.AddCommand(c)
}

// NewRootCommand creates a new root command from the given basic configuration
func NewRootCommand(c *cobra.Command) *RootCommand {
	c.TraverseChildren = true
	c.PersistentPreRunE = Setup(c.Name())
	return &RootCommand{
		c: c,
	}
}

// Context wraps a function to take a context
func Context(f func(ctx context.Context) error) func(c *cobra.Command, args []string) {
	return func(c *cobra.Command, args []string) {
		err := f(c.Context())
		if err != nil {
			log.Fatal(err)
		}
	}
}

// ContextArgs wraps a function to take a context and arguments
func ContextArgs(f func(ctx context.Context, args ...string) error) func(c *cobra.Command, args []string) {
	return func(c *cobra.Command, args []string) {
		err := f(c.Context(), args...)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// CommandContextArgs wraps a function to take a command, context and arguments
func CommandContextArgs(f func(ctx context.Context, c *cobra.Command, args ...string) error) func(c *cobra.Command, args []string) {
	return func(c *cobra.Command, args []string) {
		err := f(c.Context(), c, args...)
		if err != nil {
			log.Fatal(err)
		}
	}
}
