package types

import "fmt"

type Address struct {
	Host Host
	Port Port
}

func (a Address) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a Address) IsValid() bool {
	return a.Host.IsValid() && a.Port.IsValid()
}
