package cmd

import (
	"os"
	"path/filepath"

	"github.com/creativeprojects/filescript/term"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "filescript",
	Short: "file management on NAS",
	Long:  "\nVarious tools to help with cleaning up files on a Linux NAS (with macOS clients)",
	Run: func(cmd *cobra.Command, args []string) {
		// this function needs to be defined
	},
}

func init() {
	cobra.OnInitialize(initLog, defaultDir)
	flag := rootCmd.PersistentFlags()
	flag.StringVarP(&global.dir, "dir", "d", "", "working directory (default to current directory)")
	flag.BoolVarP(&global.write, "write", "w", false, "write updates to the disk (default to dry-run)")
	flag.BoolVarP(&global.quiet, "quiet", "q", false, "only display warnings and errors")
	flag.BoolVarP(&global.verbose, "verbose", "v", false, "display debugging information")
}

func initLog() {
	switch {
	case global.verbose:
		term.SetLevel(term.LevelDebug)
	case global.quiet:
		term.SetLevel(term.LevelWarn)
	}
	if !global.write {
		term.Info("running in dry-mode, please add '-w' flag to write to the disk")
	}
	if global.verbose {
		pterm.EnableDebugMessages()
	}
}

func defaultDir() {
	if global.dir == "" {
		global.dir, _ = os.Getwd()
	}
	global.dir = filepath.Clean(global.dir)
}

func Execute(buildVersion, buildCommit, buildDate, buildBy string) {
	term.Infof("filescript version %s built by %s (%s)", buildVersion, buildBy, buildDate)

	setApp(buildVersion, buildCommit, buildDate, buildBy)

	if err := rootCmd.Execute(); err != nil {
		term.Error(err)
		os.Exit(1)
	}
}
