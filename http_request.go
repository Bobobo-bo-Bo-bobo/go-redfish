package redfish

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func (r *Redfish) httpRequest(endpoint string, method string, header *map[string]string, reader io.Reader, basic_auth bool) (HttpResult, error) {
	var result HttpResult
	var transp *http.Transport
	var url string

	if r.InsecureSSL {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		if r.Debug {
			r.logger.Println("Enabling insecure SSL")
		}
	} else {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}

	client := &http.Client{
		Timeout:   r.Timeout,
		Transport: transp,
	}

	if r.Port > 0 {
		// check if it is an endpoint or a full URL
		if endpoint[0] == '/' {
			url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, endpoint)
		} else {
			url = endpoint
		}
	} else {
		// check if it is an endpoint or a full URL
		if endpoint[0] == '/' {
			url = fmt.Sprintf("https://%s%s", r.Hostname, endpoint)
		} else {
			url = endpoint
		}
	}

	result.Url = url

	if r.Debug {
		r.logger.Printf("Sending HTTP %s to %s with reader interface at %p\n", method, url, reader)
	}

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return result, err
	}

	if basic_auth {
		if r.Debug {
			r.logger.Printf("Setting HTTP basic authentication for HTTP %s to %s", method, url)
		}
		request.SetBasicAuth(r.Username, r.Password)
	}

	// add required headers
	request.Header.Add("OData-Version", "4.0") // Redfish API supports at least Open Data Protocol 4.0 (https://www.odata.org/documentation/)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	// set User-Agent
	request.Header.Set("User-Agent", UserAgent)

	// add authentication token if present
	if r.AuthToken != nil && *r.AuthToken != "" {
		request.Header.Add("X-Auth-Token", *r.AuthToken)
	}

	// close connection after response and prevent re-use of TCP connection because some implementations (e.g. HP iLO4)
	// don't like connection reuse and respond with EoF for the next connections
	request.Close = true

	// add supplied additional headers
	if header != nil {
		for key, value := range *header {
			request.Header.Add(key, value)
		}
	}

	if r.Debug {
		r.logger.Printf("HTTP headers for %s request to %s: %+v\n", method, url, request.Header)
	}

	response, err := client.Do(request)
	if err != nil {
		return result, err
	}

	defer response.Body.Close()

	if r.Debug {
		r.logger.Printf("HTTP %s to %s returned with status %s\n", method, url, response.Status)
		r.logger.Printf("HTTP headers returned from HTTP %s to %s: %+v\n", method, url, response.Header)
	}

	result.Status = response.Status
	result.StatusCode = response.StatusCode
	result.Header = response.Header
	result.Content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	return result, nil
}
