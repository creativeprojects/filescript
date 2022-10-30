package cmd

import (
	"context"
	"errors"

	"github.com/pterm/pterm"
)

func handleError(err error) {
	if err == nil {
		return
	}
	if errors.Is(err, context.Canceled) {
		pterm.Warning.Println("Operation cancelled")
		return
	}
	pterm.Error.Println(err)
}
