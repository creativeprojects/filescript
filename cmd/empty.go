package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/creativeprojects/filescript/fsutils"
	"github.com/creativeprojects/filescript/term"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var emptyCmd = &cobra.Command{
	Use:   "empty",
	Short: "Delete empty directories in tree",
	Run:   runEmpty,
}

func init() {
	rootCmd.AddCommand(emptyCmd)
}

func runEmpty(cmd *cobra.Command, args []string) {
	err := empty(global.dir)
	handleError(err)
}

func empty(dir string) error {
	var err error
	var totalFiles int

	if global.dir == "" {
		global.dir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	global.dir = filepath.Clean(global.dir)

	term.Debugf("searching for empty directories from %q", dir)

	spinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start()
	if err != nil {
		return err
	}

	pathCount, err := fsutils.CountFiles(context.Background(), global.dir, func(event fsutils.Event) bool {
		if event.Type == fsutils.EventProgressFile {
			totalFiles++
			spinner.Text = fsutils.Plural(totalFiles, "file")
		}
		if event.Type == fsutils.EventError {
			pterm.Error.Println(event.Err)
		}
		return true
	})
	if err != nil {
		spinner.Stop()
		return err
	}

	deletion := make([]string, 0)
	iteration := 0
	total := 0
	for {
		iteration++
		term.Debugf("\n=== New iteration %d \n", iteration)
		dirs := listEmptyDirectories(pathCount)
		if len(dirs) == 0 {
			break
		}
		for _, dir := range dirs {
			delete(pathCount, dir)
			parent := filepath.Dir(dir)
			if _, found := pathCount[parent]; found {
				pathCount[parent]--
			}
			total++
			deletion = append(deletion, dir)
			if global.write {
				err = os.Remove(dir)
				if err != nil {
					spinner.FailPrinter.Printfln("cannot delete %q: %s", dir, err)
				}
			}
		}
	}
	spinner.Stop()
	pterm.Success.Println(fsutils.Plural(total, "directory") + " deleted in " + fsutils.Plural(iteration, "iteration"))

	for _, dir := range deletion {
		term.Info(dir)
	}
	return nil
}

func listEmptyDirectories(pathCount map[string]int) (toDelete []string) {
	for path, count := range pathCount {
		term.Debugf("%s => %d\n", path, count)
		if count == 0 {
			toDelete = append(toDelete, path)
		}
	}
	return
}
