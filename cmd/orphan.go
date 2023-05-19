package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/creativeprojects/filescript/fsutils"
	"github.com/creativeprojects/filescript/term"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var orphanCmd = &cobra.Command{
	Use:     "orphan",
	Aliases: []string{"orphans"},
	Short:   "Delete macOS ._* generated files that are left alone. This does not delete the .DS_Store file.",
	Run:     runOrphan,
}

func init() {
	rootCmd.AddCommand(orphanCmd)
}

func runOrphan(cmd *cobra.Command, args []string) {
	err := orphans(global.dir)
	handleError(err)
}

func orphans(dir string) error {
	var finishedDirs, totalDirs, found int
	var err error

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	term.Debugf("searching for orphan files from %q", dir)

	spinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(true).WithText(getTextForOrphanSpinner(0, 0, 0)).Start()
	if err != nil {
		return err
	}

	progress := func(event fsutils.Event) bool {
		switch event.Type {
		case fsutils.EventError:
			spinner.Fail(event.Err, "\n")

		case fsutils.EventTotal:
			totalDirs += event.TotalDirsInDir
			spinner.Text = getTextForOrphanSpinner(finishedDirs, totalDirs, found)

		case fsutils.EventProgressDir:
			finishedDirs++
			spinner.Text = getTextForOrphanSpinner(finishedDirs, totalDirs, found)
			term.Debugf("entering %q", event.SrcDir)

		case fsutils.EventProgressFileProcessed:
			found++
			term.Debugf("orphan: %q", event.SrcFilename)
		}
		return true
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	defer stop()

	orphans, err := fsutils.FindOrphans(ctx, dir, "._", "", progress)
	_ = spinner.Stop()
	pterm.Success.Println(fsutils.Plural(found, "file") + " found")

	for _, orphan := range orphans {
		term.Info(orphan)
		if global.write {
			err = os.Remove(orphan)
			if err != nil {
				term.Error(err)
			}
		}
	}

	return err
}

func getTextForOrphanSpinner(finishedDirs, totalDirs, found int) string {
	return fmt.Sprintf("%d/%s - %s found",
		finishedDirs,
		fsutils.Plural(totalDirs, "directory"),
		fsutils.Plural(found, "file"),
	)
}
