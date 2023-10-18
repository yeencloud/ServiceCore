package error

import (
	"fmt"
	errorDomain "github.com/yeencloud/ServiceCore/src/domain/serviceError"
	"hash/crc32"
	"os"
	"runtime"
	"strings"
)

func stackFingerprint(s errorDomain.Stack) string {
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

func getStack() errorDomain.Stack {
	stack := make(errorDomain.Stack, 0)

	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		functionName := functionName(pc)

		stack = append(stack, errorDomain.Frame{file, functionName, line})
	}

	//stack = stack[:len(stack)-1]

	stack = stack[2 : len(stack)-1]

	return stack
}

func fingerprint(str string) string {
	hash := crc32.NewIEEE()

	fmt.Fprintf(hash, str)

	return fmt.Sprintf("%x", hash.Sum32())
}

func fingerprintError(trace *errorDomain.Error) string {
	hash := ""
	currentTrace := trace
	for {
		hash = fmt.Sprintf("%s%s", hash, stackFingerprint(currentTrace.Stack))
		currentTrace = currentTrace.Child
		if currentTrace == nil {
			break
		}
	}
	return fingerprint(hash)
}

func Trace(err errorDomain.ErrorDescription) *errorDomain.Error {
	error := errorDomain.Error{
		ErrorDescription: err,
		Stack:            getStack(),
	}

	error.Fingerprint = fingerprintError(&error)

	return &error
}
