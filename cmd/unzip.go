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
		deletion := []string{}
		for filename := range eventChan {
			if global.write {
				pterm.Debug.Println(filename)
				err := fsutils.Unzip(context.Background(), filename, []string{"thumbs"}, progress)
				if err != nil {
					pterm.Error.Println(err)
					continue
				}
				deletion = append(deletion, filename)
			} else {
				pterm.Info.Printf("would unzip %q\n", filename)
			}
		}
		pterm.Success.Printf("deleting %d files\n", len(deletion))
		for _, filename := range deletion {
			err := os.Remove(filename)
			if err != nil {
				pterm.Error.Println(err)
			}
		}
		wg.Done()
	}()

	err = fsutils.FindFiles(context.Background(), fsutils.WithExtension(".zip"), dir, eventChan, progress)
	close(eventChan)
	wg.Wait()
	time.Sleep(2 * time.Second)
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
