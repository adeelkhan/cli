package main

import (
	"fmt"

	"github.com/docker/cli/cli-plugins/metadata"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func main() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		return &cobra.Command{
			Use:   "nopersistentprerun",
			Short: "Testing without PersistentPreRun hooks",
			// PersistentPreRunE: Not specified, we need to test that it works in the absence of an explicit call
			RunE: func(cmd *cobra.Command, args []string) error {
				cli := dockerCli.Client()
				ping, err := cli.Ping(cmd.Context())
				if err != nil {
					return err
				}
				fmt.Println(ping.APIVersion)
				return nil
			},
		}
	},
		metadata.Metadata{
			SchemaVersion: "0.1.0",
			Vendor:        "Docker Inc.",
			Version:       "testing",
		})
}
