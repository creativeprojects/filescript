package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/creativeprojects/filescript/fsutils"
	"github.com/creativeprojects/filescript/term"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var dirModCmd = &cobra.Command{
	Use:   "dirmod",
	Short: "Add execution permission to directories.",
	Run:   runDirMod,
}

func init() {
	rootCmd.AddCommand(dirModCmd)
}

func runDirMod(cmd *cobra.Command, args []string) {
	err := dirMod(global.dir)
	handleError(err)
}

func dirMod(dir string) error {
	var err error
	var totalDirs, found int

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	term.Debugf("fixing dir mod from %q", dir)

	spinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(true).WithText(getTextForDirModSpinner(0, 0)).Start()
	if err != nil {
		return err
	}

	progress := func(event fsutils.Event) bool {
		switch event.Type {
		case fsutils.EventError:
			pterm.Error.Println(event.Err)

		case fsutils.EventProgressDir:
			totalDirs++
			spinner.Text = getTextForDirModSpinner(totalDirs, found)

		}
		return true
	}

	var eventChan = make(chan string, 1000)
	wg := sync.WaitGroup{}
	wg.Go(func() {
		for filename := range eventChan {
			info, err := os.Stat(filename)
			if err != nil {
				pterm.Error.Println(err)
				continue
			}
			mode := info.Mode() | 0o111
			if global.write {
				err := os.Chmod(filename, mode)
				if err != nil {
					pterm.Error.Println(err)
					continue
				}
			} else {
				pterm.Info.Printf("would fix %q from %#o to %#o\n", filename, info.Mode(), mode)
			}
		}
	})

	err = fsutils.FindDirs(context.Background(), fsutils.WithoutExecutionBit(), dir, eventChan, progress)
	close(eventChan)
	wg.Wait()
	time.Sleep(2 * time.Second)
	_ = spinner.Stop()
	return err
}

func getTextForDirModSpinner(totalDirs, found int) string {
	return fmt.Sprintf("%s found",
		fsutils.Plural(totalDirs, "directory"),
	)
}
