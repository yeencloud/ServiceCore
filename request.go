package servicecore

type Version int

const APIVersion = Version(1)

type ServiceRequest struct {
	Service    string
	Method     string
	ApiVersion Version

	Request map[string]interface{}
}