package cmds

import "errors"

var (
	ErrNotHaveProcess      = errors.New(`process don't exist'`)
	ErrNotFoundServiceDir  = errors.New(`service dir don't exist'`)
	ErrNotFoundServerIndex = errors.New(`server index don't found'`)
)

const (
	ErrTextOpenFailed  = `open %s %v`
	ErrTextCloseFailed = `close %s %v`
	ErrTextReadFailed  = `read %s %v`
	ErrTextSeekFailed  = `%s seek %v`
	ErrTextWriteFailed = `write to %s %v`
)

// text 内容
const (
	TextNameResult       = "%s %v -> %v"
	TextNameCloseSuccess = "close %s success"
)

// log 内容
const (
	LogNameServiceRun = "run cmd: %s, args: %v, dir: %v\n"
	LogNameSetPort    = "set port %s = %d"
)
