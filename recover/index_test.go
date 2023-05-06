package recover

import "testing"

func TestRecoverFromPanic(t *testing.T) {
	defer RecoverFromPanic()

	panic("test panic")
}
