package galaxy

type LookUpRequest struct {
	Service string
	Method  string
}

type LookUpResponse struct {
	Address string
	Port    int
}
