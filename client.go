package ServiceCore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
	"os"
)

func (sh *ServiceHost) Call(service string, method string, args any) (map[string]interface{}, error) {
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

	requestURL := fmt.Sprintf("http://localhost:%d/rpc/", serverPort)

	b, _ := json.Marshal(&args)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	j, _ := json.Marshal(ServiceRequest{
		Service: service,
		Method:  method,
		Version: 1,
		Request: m,
	})

	spew.Dump(string(j))

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(j))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	return response, nil
}