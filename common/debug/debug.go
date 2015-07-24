package debug

import (
	"fmt"
	"os"
)

type (
	DebugLevelT int
)

const (
	DEBUG0 DebugLevelT = iota
	DEBUG1
	DEBUG2
	DEBUG3
)

// Debugging can be used to check if debugging is enabled
var (
	DebugLevel = DEBUG0
	levelColor = map[DebugLevelT]string{
		DEBUG0: "",
		DEBUG1: "\033[38;5;244m",
		DEBUG2: "\033[38;5;239m",
		DEBUG3: "\033[38;5;235m",
	}
)

// Debug contains function to print debug messages. Defaults to DebugNone (nothing printed)
var Debug = func(lvl DebugLevelT, msg string, args ...interface{}) {
	if lvl > DebugLevel {
		return
	}
	l := len(msg)
	if l > 1 && msg[l-1] != '\n' {
		msg += "\n"
	}
	msg = fmt.Sprintf("%sDBG%d:\033[0m %s", levelColor[lvl], lvl, msg)
	os.Stderr.Write([]byte(fmt.Sprintf(msg, args...)))
}
