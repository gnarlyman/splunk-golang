package splunk

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// RestartServer ...
func (conn SplunkConnection) RestartServer() (string, error) {
	response, err := conn.httpPost(fmt.Sprintf("%s/services/server/control/restart", conn.BaseURL), nil)
	return response, err
}

// EnableHEC ...
func (conn SplunkConnection) EnableHEC(data *url.Values) (string, error) {
	response, err := conn.httpPost(fmt.Sprintf("%s/servicesNS/admin/splunk_httpinput/data/inputs/http/http/enable", conn.BaseURL), data)
	return response, err
}

// CreateHECEndpoint ...
func (conn SplunkConnection) CreateHECEndpoint(data *url.Values) (CreateHECEndpointResponse, error) {
	response, err := conn.httpPost(fmt.Sprintf("%s/services/data/inputs/http", conn.BaseURL), data)
	r := []byte(response)
	var d CreateHECEndpointResponse
	json.Unmarshal(r, &d)
	return d, err
}

// DeleteHECEndpoint ...
func (conn SplunkConnection) DeleteHECEndpoint(data *url.Values, name string) (string, error) {
	var err error
	response, err := conn.httpDelete(fmt.Sprintf("%s/services/data/inputs/http/%s", conn.BaseURL, name), data)
	return response, err
}

// CreateHECEndpointResponse ...
type CreateHECEndpointResponse struct {
	Entry []struct {
		Name    string `json:"name"`
		Content struct {
			Token string `json:"token"`
		} `json:"content"`
	} `json:"entry"`
}
