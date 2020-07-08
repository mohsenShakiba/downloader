package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	downloadUrl string
)

func GetAddCommand() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "add a link for download, app will begin downloading immediately",
	}

	addCmd.Flags().StringVarP(&downloadUrl, "url", "u", "", "the url of content to be downloaded")
	err := addCmd.MarkFlagRequired("url")

	if err != nil {
		fmt.Printf("failed to run, error: %s", err)
		os.Exit(1)
	}

	addCmd.Run = func(cmd *cobra.Command, args []string) {
		downloadManager.AddDownload(downloadUrl)
	}

	return addCmd
}
