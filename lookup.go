package ServiceCore

type LookUpRequest struct {
	ServiceMethod string
}

type LookUpResponse struct {
	Address string
}

func (galaxy *GalaxyClient) LookUp(serviceMethod string) (string, error) {
	var response LookUpResponse
	err := galaxy.client.Call("Galaxy.LookUp", LookUpRequest{serviceMethod}, &response)
	if err != nil {
		return "", err
	}
	return response.Address, nil
}
