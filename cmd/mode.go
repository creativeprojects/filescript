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

var fileModeCmd = &cobra.Command{
	Use:     "mode",
	Aliases: []string{"filemode", "unexecute"},
	Short:   "Remove execution permission from files, add execution permission to directories.",
	Run:     runFileMode,
}

func init() {
	rootCmd.AddCommand(fileModeCmd)
}

func runFileMode(cmd *cobra.Command, args []string) {
	err := filemode(global.dir)
	handleError(err)
}

func filemode(dir string) error {
	var err error
	var totalDirs, found int

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	term.Debugf("fixing file mode from %q", dir)

	spinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(true).WithText(getTextForFilemodeSpinner(0, 0)).Start()
	if err != nil {
		return err
	}

	progress := func(event fsutils.Event) bool {
		switch event.Type {
		case fsutils.EventError:
			pterm.Error.Println(event.Err)

		case fsutils.EventProgressDir:
			totalDirs++
			spinner.Text = getTextForFilemodeSpinner(totalDirs, found)

		case fsutils.EventProgressFile:
			found++
			spinner.Text = getTextForFilemodeSpinner(totalDirs, found)
		}
		return true
	}

	var eventChan = make(chan string, 1000)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for filename := range eventChan {
			info, err := os.Stat(filename)
			if err != nil {
				pterm.Error.Println(err)
				continue
			}
			mode := info.Mode() & 0o666
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
		wg.Done()
	}()

	err = fsutils.FindFiles(context.Background(), fsutils.WithExecutionBit(), dir, eventChan, progress)
	close(eventChan)
	wg.Wait()
	time.Sleep(2 * time.Second)
	_ = spinner.Stop()
	return err
}

func getTextForFilemodeSpinner(totalDirs, found int) string {
	return fmt.Sprintf("%s - %s found",
		fsutils.Plural(totalDirs, "directory"),
		fsutils.Plural(found, "file"),
	)
}
