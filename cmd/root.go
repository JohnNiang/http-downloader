package cmd

import (
	"context"
	extver "github.com/linuxsuren/cobra-extension/version"
	"github.com/spf13/cobra"
)

// NewRoot returns the root command
func NewRoot(cxt context.Context) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "hd",
		Short: "HTTP download tool",
	}

	cmd.AddCommand(
		NewGetCmd(cxt), NewInstallCmd(cxt), newFetchCmd(cxt), newSearchCmd(cxt),
		extver.NewVersionCmd("linuxsuren", "http-downloader", "hd", nil))
	return
}
