package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

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
	var totalDirs, found int

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	term.Debugf("searching zip files from %q", dir)

	spinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(true).WithText(getTextForUnzipSpinner(0, 0)).Start()
	if err != nil {
		return err
	}

	progress := func(event fsutils.Event) bool {
		switch event.Type {
		case fsutils.EventError:
			pterm.Error.Println(event.Err)

		case fsutils.EventProgressDir:
			totalDirs++
			spinner.Text = getTextForUnzipSpinner(totalDirs, found)

		case fsutils.EventProgressFile:
			found++
		}
		return true
	}

	var eventChan = make(chan string, 1000)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for filename := range eventChan {
			pterm.Debug.Println(filename)
			dir := filepath.Join(filepath.Dir(filename), filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))]))
			if exists, err := fsutils.Exists(dir); exists && err == nil {
				pterm.Warning.Printfln("directory %q already exists, skipping", dir)
			}
		}
		wg.Done()
	}()

	err = fsutils.FindFiles(context.Background(), fsutils.WithExtension(".zip"), dir, eventChan, progress)
	close(eventChan)
	wg.Wait()
	spinner.Stop()
	return err
}

func getTextForUnzipSpinner(totalDirs, found int) string {
	return fmt.Sprintf("%s - %s found",
		fsutils.Plural(totalDirs, "directory"),
		fsutils.Plural(found, "archive file"),
	)
}
