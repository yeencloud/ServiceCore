package servicecore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yeencloud/ServiceCore/domain"
	"github.com/yeencloud/ServiceCore/tools"
	"io/ioutil"
	"net/http"
	"os"
)

func (sh *ServiceHost) callWithAddress(hostname string, port int, service string, method string, args any) (map[string]interface{}, *domain.ServiceError) {
	requestURL := fmt.Sprintf("http://%s:%d/rpc/", hostname, port)

	callData := tools.AnyToMap(args)

	j, _ := json.Marshal(domain.ServiceRequest{
		Module:     service,
		Method:     method,
		ApiVersion: domain.APIVersion,
		Data:       callData,
	})

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(j))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, ErrCallCouldNotReadResponseBody.Embed(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, ErrCallCouldNotReadResponseBody.Embed(err)
	}

	var response map[string]interface{}
	err = json.Unmarshal(resBody, &response)

	if err != nil {
		return nil, ErrCallCouldNotUnmarshalResponseBody.Embed(err)
	}

	return response, nil
}

// Call calls a service method with the given arguments (preferably a struct).
func (sh *ServiceHost) Call(service string, method string, args any) (map[string]interface{}, *domain.ServiceError) {
	address, err := sh.LookUp(service, method)
	if err != nil {
		return nil, err
	}

	serverPort := 8000

	return sh.callWithAddress(address, serverPort, service, method, args)
}