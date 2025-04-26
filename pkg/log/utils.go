package logger

import (
	"fmt"
	"runtime"
	"strings"
)

func getFullStack() string {
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, true)
	stack := fmt.Sprintf("%s", buf[0:stackSize])
	stackTemp := strings.Split(stack, "\n")
	stackFile := fmt.Sprintf("file: %s, func: %s", strings.TrimSpace(stackTemp[6]), strings.TrimSpace(stackTemp[5]))
	return stackFile
}
