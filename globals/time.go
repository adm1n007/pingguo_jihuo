package globals

import (
	"time"
)

func GetCurrentTime() int64 {
	// milliseconds
	return time.Now().UnixNano() / int64(time.Millisecond)
}
