package cmd

type GlobalFlags struct {
	dir     string
	quiet   bool
	verbose bool
}

var (
	global GlobalFlags
)
