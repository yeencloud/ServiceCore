package ServiceCore

type CoreHealth struct {
}

type HealthCall struct {
}

type HealthResponse struct {
	Status bool
	Msg    string
}

func (ch *CoreHealth) Health(args *HealthCall, response *HealthResponse) error {
	response.Status = true
	response.Msg = "unimplemented at the moment"

	return nil
}