package cmd

import (
	"github.com/linuxsuren/http-downloader/pkg/exec"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"
	"path"
)

func newFetchCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "fetch",
		Short: "Fetch the latest hd config",
		RunE: func(_ *cobra.Command, _ []string) (err error) {
			return fetchHomeConfig()
		},
	}
	return
}

func getConfigDir() (configDir string, err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err == nil {
		configDir = path.Join(userHome, "/.config/hd-home")
	}
	return
}

func fetchHomeConfig() (err error) {
	var configDir string
	if configDir, err = getConfigDir(); err != nil {
		return
	}

	if ok, _ := pathExists(configDir); ok {
		err = exec.ExecCommandInDir("git", configDir, "reset", "--hard", "origin/master")
		if err == nil {
			err = exec.ExecCommandInDir("git", configDir, "pull")
		}
	} else {
		if err = os.MkdirAll(configDir, 0644); err == nil {
			err = exec.ExecCommand("git", "clone", "https://github.com/LinuxSuRen/hd-home", configDir)
		}
	}
	return
}
