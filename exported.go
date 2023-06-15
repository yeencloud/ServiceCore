package ServiceCore

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go/token"
	"reflect"
	"sync"
)

func ExportList(rcvr any) []string {
	server := new(Server)

	server.register(rcvr, "", false)

	exportList := []string{}

	server.serviceMap.Range(func(key, value interface{}) bool {
		service := value.(*Service)
		name := key.(string)

		for k, _ := range service.Method {
			value := k
			exportList = append(exportList, fmt.Sprintf("%s.%s", name, value))
		}

		return true
	})
	return exportList
}

// ------------------ copied from net/rpc/server.go ------------------
// golang doesn't export those functions, so we have to copy them here to be able to find all the exported methods for the rpc registration
type Service struct {
	Name   string                 // name of service
	Rcvr   reflect.Value          // receiver of methods for the service
	Typ    reflect.Type           // type of the receiver
	Method map[string]*MethodType // registered methods
}

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type MethodType struct {
	Method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

// Server represents an RPC Server.
type Server struct {
	serviceMap sync.Map // map[string]*service
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return token.IsExported(t.Name()) || t.PkgPath() == ""
}

func (server *Server) register(rcvr any, name string, useName bool) error {
	s := new(Service)
	s.Typ = reflect.TypeOf(rcvr)
	s.Rcvr = reflect.ValueOf(rcvr)
	sname := name
	if !useName {
		sname = reflect.Indirect(s.Rcvr).Type().Name()
	}
	if sname == "" {
		log.Err(errNoServiceNameForType).Str("type", s.Typ.String())
		return errNoServiceNameForType
	}
	if !useName && !token.IsExported(sname) {
		log.Err(errTypeIsNotExported).Str("type", sname)
		return errTypeIsNotExported
	}
	s.Name = sname

	// Install the methods
	s.Method = suitableMethods(s.Typ)

	if len(s.Method) == 0 {
		method := suitableMethods(reflect.PointerTo(s.Typ))
		log.Err(errTypeHasNoSuitableExportedMethods).Int("method", len(s.Method)).Int("pointedMethods", len(method)).Str("type", sname)
		return errTypeHasNoSuitableExportedMethods
	}

	if _, dup := server.serviceMap.LoadOrStore(sname, s); dup {
		log.Err(errServiceAlreadyDefined).Int("method", len(s.Method)).Str("service", sname)
		return errServiceAlreadyDefined
	}
	return nil
}

// suitableMethods returns suitable Rpc methods of typ. It will log
// errors if logErr is true.
func suitableMethods(typ reflect.Type) map[string]*MethodType {
	methods := make(map[string]*MethodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if !method.IsExported() {
			continue
		}
		// Method needs three ins: receiver, *args, *reply.
		if mtype.NumIn() != 3 {
			log.Warn().Str("method", mname).Int("in", mtype.NumIn()).Msg("method has wrong number of ins, should be 3")
			continue
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			log.Warn().Str("method", mname).Str("argType", argType.String()).Msg("argument type not exported")
			continue
		}
		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Pointer {
			log.Warn().Str("method", mname).Str("replyType", replyType.String()).Msg("reply type not a pointer")
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			log.Warn().Str("method", mname).Str("replyType", replyType.String()).Msg("reply type not exported")
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			log.Warn().Str("method", mname).Int("out", mtype.NumOut()).Msg("method has wrong number of outs, should be 1")
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			log.Warn().Str("method", mname).Str("returnType", returnType.String()).Msg("return type not error")
			continue
		}
		methods[mname] = &MethodType{Method: method, ArgType: argType, ReplyType: replyType}
	}
	return methods
}