package cmd

type GlobalFlags struct {
	dir     string
	write   bool
	quiet   bool
	verbose bool
}

var (
	global GlobalFlags
)
