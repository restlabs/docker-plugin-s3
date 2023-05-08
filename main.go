package main

import (
	"encoding/json"
	"fmt"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/roshbhatia/docker-plugin-s3/pkg/commands"
	"github.com/spf13/cobra"
)

var metadata = manager.Metadata{
	SchemaVersion:     "1.0.0",
	Vendor:            "roshbhatia",
	Version:           "1.0.0",
	ShortDescription:  "Upload and download Docker images to S3 compatible storage",
	Experimental:      false,
}

func main() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "s3",
			Short: "Docker CLI plugin to upload Docker images to S3 compatible storage",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				if err := plugin.PersistentPreRunE(cmd, args); err != nil {
					return err
				}
				return nil
			},
		}

		cmd.AddCommand(&cobra.Command{
			Use:   "docker-cli-plugin-metadata",
			Short: "Print plugin metadata",
			Run: func(cmd *cobra.Command, _ []string) {
				metadataBytes, err := json.Marshal(metadata)
				if err != nil {
					fmt.Fprintln(dockerCli.Err(), err)
					return
				}
				fmt.Fprintln(dockerCli.Out(), string(metadataBytes))
				return
			},
		})

		cmd.AddCommand(commands.PullCmd)
		cmd.AddCommand(commands.PushCmd)

		return cmd
	}, metadata)
}
