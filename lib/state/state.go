package state

type State int

const (
	OK State = iota
	WARN
	CRITICAL
	UNKNOWN
)
