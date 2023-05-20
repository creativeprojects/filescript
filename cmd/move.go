package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/creativeprojects/filescript/fsutils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var runMoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move files in subfolders per year",
	Run:   runMove,
}

func init() {
	rootCmd.AddCommand(runMoveCmd)
}

func runMove(cmd *cobra.Command, args []string) {
	err := movePerYear(global.dir)
	handleError(err)
}

func movePerYear(dir string) error {
	var err error
	var filesFound, filesProcessed int

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	pterm.Debug.Printf("searching zip files from %q\n", dir)

	spinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(true).WithText(getTextForMovePerYear(0, 0)).Start()
	if err != nil {
		return err
	}

	progress := func(event fsutils.Event) bool {
		switch event.Type {
		case fsutils.EventError:
			pterm.Error.Println(event.Err)

		case fsutils.EventProgressFile:
			filesFound++
			spinner.Text = getTextForMovePerYear(filesFound, filesProcessed)

		case fsutils.EventProgressFileProcessed:
			filesProcessed++
			spinner.Text = getTextForMovePerYear(filesFound, filesProcessed)
		}
		return true
	}

	err = fsutils.MoveAllPerYear(context.Background(), dir, progress, !global.write)
	_ = spinner.Stop()
	return err
}

func getTextForMovePerYear(filesFound, filesProcessed int) string {
	return fmt.Sprintf("%s found - %s moved",
		fsutils.Plural(filesFound, "file"),
		fsutils.Plural(filesProcessed, "file"),
	)
}
