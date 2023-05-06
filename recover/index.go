package recover

import (
	"log"
	"runtime"
)

// RecoverFromPanic 从panic中恢复,并且记录错误的位置和错误信息
func RecoverFromPanic() {
	if err := recover(); err != nil {
		buf := make([]byte, 64<<10)
		buf = buf[:runtime.Stack(buf, false)]
		log.Printf("panic recover err: %s\n%s", err, buf)
	}
}
