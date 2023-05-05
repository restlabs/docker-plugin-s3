package main

import (
	"encoding/json"
	"fmt"
	"os"

	upload "github.com/roshbhatia/docker-plugin-s3/pkg/cmd"
	"github.com/spf13/cobra"
)

type PluginMetadata struct {
	SchemaVersion    string   `json:"SchemaVersion"`
	Vendor           string   `json:"Vendor"`
	Name             string   `json:"Name"`
	Version          string   `json:"Version"`
	ShortDescription string   `json:"ShortDescription"`
	DockerVersion    string   `json:"DockerVersion"`
	Experimental     bool     `json:"Experimental"`
	Platforms        []string `json:"Platforms"`
}

var metadata = PluginMetadata{
	SchemaVersion:    "0.1.0",
	Vendor:           "Rosh Bhatia",
	Name:             "s3",
	Version:          "0.1.0",
	ShortDescription: "Upload Docker images to S3 compatible storage",
	DockerVersion:    ">=20.10.0",
	Experimental:     false,
	Platforms:        []string{"linux/amd64", "darwin/amd64", "windows/amd64"},
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "docker-cli-plugin-metadata" {
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to generate metadata JSON")
			os.Exit(1)
		}
		fmt.Println(string(metadataJSON))
		os.Exit(0)
	}

	cmd := &cobra.Command{
		Use:              "s3",
		Short:            "A Docker CLI plugin to upload Docker images to S3 compatible storage",
		TraverseChildren: true,
	}	
	

	cmd.AddCommand(upload.UploadCmd)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
