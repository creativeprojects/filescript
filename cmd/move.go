package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/creativeprojects/filescript/fsutils"
	"github.com/spf13/cobra"
)

var runMoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move files in subfolders per year",
	RunE:  runMove,
}

func init() {
	rootCmd.AddCommand(runMoveCmd)
}

func runMove(cmd *cobra.Command, args []string) error {
	var dir string
	var err error

	flag.StringVar(&dir, "d", "", "directory where to move files per year")
	flag.Parse()

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("cannot stat %q: %s", entry.Name(), err)
			continue
		}
		year := strconv.Itoa(info.ModTime().Year())
		if len(year) != 4 {
			continue
		}
		// fmt.Printf("%s: %q\n", year, entry.Name())
		err = os.MkdirAll(filepath.Join(dir, year), 0777)
		if err != nil {
			log.Fatal(err)
		}
		orig := filepath.Join(dir, entry.Name())
		moveTo := filepath.Join(dir, year, entry.Name())

		newpath, err := fsutils.Rename(orig, moveTo)
		if err != nil {
			log.Fatal(err)
		}
		if newpath != moveTo {
			fmt.Printf("* file \n   %q renamed to\n   %q\n", orig, newpath)
		}
	}
	return nil
}
