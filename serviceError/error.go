package serviceError

import (
	"fmt"
)

type ErrorDescription struct {
	HttpCode int
	String   string
}

type Error struct {
	ErrorDescription

	Stack           Stack
	Child           *Error                 `json:",omitempty"`
	AdditionnalData map[string]interface{} `json:",omitempty"`
	Fingerprint     string
}

func Trace(err ErrorDescription) *Error {
	error := Error{
		ErrorDescription: err,
		Stack:            getStack(),
	}

	error.Fingerprint = error.fingerprint()

	return &error
}

func Wrap(err error) *Error {
	/*if err == nil {
		return nil
	}
	return Trace(err)*/
	return nil
}

func (trace *Error) Append(err error) *Error {
	/*if trace == nil || err == nil {
		return nil
	}

	newErr := Trace(err)
	newErr.Child = trace

	return newErr*/

	return nil
}

func (trace *Error) fingerprint() string {
	hash := ""
	currentTrace := trace
	for {
		hash = fmt.Sprintf("%s%s", hash, currentTrace.Stack.Fingerprint())
		currentTrace = currentTrace.Child
		if currentTrace == nil {
			break
		}
	}
	return fingerprint(hash)
}