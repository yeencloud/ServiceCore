package types

type Port int

func (p Port) IsValid() bool {
	return p > 0 && p < 65535
}
