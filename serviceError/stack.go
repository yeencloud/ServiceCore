package serviceError

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type Stack []Frame

func (s Stack) Fingerprint() string {
	hash := ""
	for _, frame := range s {
		hash = fmt.Sprintf("%s%s%s%d", hash, frame.File, frame.Method, frame.Line)
	}
	return fingerprint(hash)
}

func functionName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "<unknown>"
	}
	name := fn.Name()
	end := strings.LastIndex(name, string(os.PathSeparator))
	return name[end+1:]
}

func getStack() Stack {
	stack := make(Stack, 0)

	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		functionName := functionName(pc)

		if !strings.HasPrefix(functionName, "go-nested-traced-error") {
			stack = append(stack, Frame{file, functionName, line})
		}
	}

	//stack = stack[:len(stack)-1]

	stack = stack[2 : len(stack)-1]

	return stack
}