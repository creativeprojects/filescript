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

var runUnzipCmd = &cobra.Command{
	Use:   "unzip",
	Short: "Unzip all archive files in subfolders",
	Run:   runUnzip,
}

func init() {
	rootCmd.AddCommand(runUnzipCmd)
}

func runUnzip(cmd *cobra.Command, args []string) {
	err := unzip(global.dir)
	handleError(err)
}

func unzip(dir string) error {
	var err error
	var totalDirs, found, unzipped int

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	term.Debugf("searching zip files from %q", dir)

	spinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(true).WithText(getTextForUnzipSpinner(0, 0, 0)).Start()
	if err != nil {
		return err
	}

	progress := func(event fsutils.Event) bool {
		switch event.Type {
		case fsutils.EventError:
			pterm.Error.Println(event.Err)

		case fsutils.EventProgressDir:
			totalDirs++
			spinner.Text = getTextForUnzipSpinner(totalDirs, found, unzipped)

		case fsutils.EventProgressFile:
			found++
			spinner.Text = getTextForUnzipSpinner(totalDirs, found, unzipped)

		case fsutils.EventProgressFileProcessed:
			unzipped++
			spinner.Text = getTextForUnzipSpinner(totalDirs, found, unzipped)
		}
		return true
	}

	var eventChan = make(chan string, 1000)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for filename := range eventChan {
			pterm.Debug.Println(filename)
			err := fsutils.Unzip(context.Background(), filename, nil, progress)
			if err != nil {
				pterm.Error.Println(err)
			}
			time.Sleep(5 * time.Second)
		}
		wg.Done()
	}()

	err = fsutils.FindFiles(context.Background(), fsutils.WithExtension(".zip"), dir, eventChan, progress)
	close(eventChan)
	wg.Wait()
	_ = spinner.Stop()
	return err
}

func getTextForUnzipSpinner(totalDirs, found, unzipped int) string {
	return fmt.Sprintf("%s - %s found - %s unzipped",
		fsutils.Plural(totalDirs, "directory"),
		fsutils.Plural(found, "archive file"),
		fsutils.Plural(unzipped, "file"),
	)
}
