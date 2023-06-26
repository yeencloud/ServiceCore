package servicecore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
)

func (sh *ServiceHost) callWithAddress(hostname string, port int, service string, method string, args any) (map[string]interface{}, error) {

	requestURL := fmt.Sprintf("http://%s:%d/rpc/", hostname, port)

	b, _ := json.Marshal(&args)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	j, _ := json.Marshal(ServiceRequest{
		Service:    service,
		Method:     method,
		ApiVersion: APIVersion,
		Request:    m,
	})

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(j))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Err(err).Msg("could not make http request")
		os.Exit(1)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Err(err).Msg("could not read response body")
		os.Exit(1)
	}

	var response map[string]interface{}
	json.Unmarshal(resBody, &response)

	return response, nil
}

// Call calls a service method with the given arguments (preferably a struct).
func (sh *ServiceHost) Call(service string, method string, args any) (map[string]interface{}, error) {
	address, err := sh.LookUp(service, method)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(args)
	if err != nil {
		return map[string]interface{}{}, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(marshal, &response)
	if err != nil {
		return map[string]interface{}{}, err
	}

	serverPort := 8000

	return sh.callWithAddress(address, serverPort, service, method, args)
}