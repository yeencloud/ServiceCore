package ServiceCore

type ApiVersion int

type ServiceRequest struct {
	Service string
	Method  string
	Version ApiVersion

	Request map[string]interface{}
}