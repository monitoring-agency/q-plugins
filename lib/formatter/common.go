package formatter

type State int

const (
	OK State = iota + 1
	WARN
	CRITICAL
	UNKNOWN
)
