package domain

type Version int

const APIVersion = Version(1)

type ServiceRequest struct {
	Module     string
	Method     string
	ApiVersion Version

	Data map[string]interface{}
}