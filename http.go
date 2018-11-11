package splunk

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*
 * HTTP helper methods
 */

func httpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}

func (conn SplunkConnection) httpGet(url string, data *url.Values) (string, error) {
	var response *http.Response
	var err error
	if response, err = conn.httpCall(url, "GET", data); err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	return string(body), nil
}

func (conn SplunkConnection) httpPost(url string, data *url.Values) (string, error) {
	var response *http.Response
	var err error
	if response, err = conn.httpCall(url, "POST", data); err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	return string(body), nil
}

func (conn SplunkConnection) httpDelete(url string, data *url.Values) (string, error) {
	var response *http.Response
	var err error
	if response, err = conn.httpCall(url, "DELETE", data); err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	return string(body), nil
}

func (conn SplunkConnection) httpCall(url string, method string, data *url.Values) (*http.Response, error) {
	client := httpClient()

	var payload io.Reader
	if data != nil {
		payload = bytes.NewBufferString(data.Encode())
	}

	request, err := http.NewRequest(method, url, payload)
	conn.addAuthHeader(request)
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	return response, err
}

func (conn SplunkConnection) addAuthHeader(request *http.Request) {
	if conn.sessionKey.Value != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Splunk %s", conn.sessionKey))
	} else {
		request.SetBasicAuth(conn.Username, conn.Password)
	}
}
