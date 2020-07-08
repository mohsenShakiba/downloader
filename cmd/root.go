package cmd

import (
	"downloader/manager"
	"github.com/spf13/cobra"
	"os"
)

var downloadManager *manager.Manager

func init() {
	downloadManager = manager.NewManager(os.TempDir())
}

func Execute() error {

	rootCmd := &cobra.Command{
		Use:   "downloaded",
		Short: "A micro app for downloading links efficiently from internet",
	}

	addCmd := GetAddCommand()

	rootCmd.AddCommand(addCmd)

	return rootCmd.Execute()
}
