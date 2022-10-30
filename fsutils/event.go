package fsutils

type EventType int

const (
	EventError EventType = iota
	EventTotal
	EventProgressFile
	EventProgressDir
	EventProgressFileProcessed
	EventProgressDirProcessed
)

type Event struct {
	Type            EventType
	Err             error  // only available when event Type is EventError
	TotalFilesInDir int    // only available when event Type is EventTotal
	TotalDirsInDir  int    // only available when event Type is EventTotal
	SrcDir          string // only available when event Type is EventProgressFileProcessed or EventProgressDirProcessed
	SrcFilename     string // only available when event Type is EventProgressFileProcessed or EventProgressDirProcessed
	DstDir          string // only available when event Type is EventProgressFileProcessed or EventProgressDirProcessed
	DstFilename     string // only available when event Type is EventProgressFileProcessed or EventProgressDirProcessed
}
