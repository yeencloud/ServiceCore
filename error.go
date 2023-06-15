package ServiceCore

import "errors"

var (
	errNoServiceNameForType             = errors.New("no service name for type")
	errTypeIsNotExported                = errors.New("type is not exported")
	errTypeHasNoSuitableExportedMethods = errors.New("type has no suitable exported methods")
	errServiceAlreadyDefined            = errors.New("service already defined")
)