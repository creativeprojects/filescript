package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/creativeprojects/filescript/fsutils"
	"github.com/spf13/cobra"
)

var emptyCmd = &cobra.Command{
	Use:   "empty",
	Short: "Delete empty directories in tree",
	RunE:  runEmpty,
}

func init() {
	rootCmd.AddCommand(emptyCmd)
}

func runEmpty(cmd *cobra.Command, args []string) error {
	var dir string
	var remove bool
	var err error

	flag.StringVar(&dir, "d", "", "directory to remove when empty")
	flag.BoolVar(&remove, "r", false, "remove empty directories")
	flag.Parse()

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	dir = filepath.Clean(dir)

	pathCount, err := fsutils.CountFiles(context.Background(), dir, func(event fsutils.Event) bool {
		return true
	})
	if err != nil {
		log.Fatal(err)
	}

	iteration := 0
	total := 0
	for {
		iteration++
		fmt.Printf("\n=== New iteration %d \n", iteration)
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
			if remove {
				err = os.Remove(dir)
				if err != nil {
					log.Fatalf("cannot delete %q: %s", dir, err)
				}
			}
		}
	}
	fmt.Printf("\n Total %d dir deleted in %d iterations\n", total, iteration)
	return nil
}

func listEmptyDirectories(pathCount map[string]int) (toDelete []string) {
	for path, count := range pathCount {
		fmt.Printf("%s => %d\n", path, count)
		if count == 0 {
			toDelete = append(toDelete, path)
		}
	}
	return
}
