package redfish

import (
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

func (r *Redfish) httpRequest(endpoint string, method string, header *map[string]string, reader io.Reader, basicAuth bool) (HTTPResult, error) {
	var result HTTPResult
	var transp *http.Transport
	var url string

	if r.InsecureSSL {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		if r.Debug {
			log.Debug("Disabling verification of SSL certificates")
		}
	} else {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}

	client := &http.Client{
		Timeout:   r.Timeout,
		Transport: transp,
		// non-GET methods (like PATCH, POST, ...) may or may not work when encountering
		// HTTP redirect. Don't follow 301/302. The new location can be checked by looking
		// at the "Location" header.
		CheckRedirect: func(http_request *http.Request, http_via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	if r.Port > 0 && r.Port != 443 {
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

	result.URL = url

	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":      r.Hostname,
			"port":          r.Port,
			"timeout":       r.Timeout,
			"flavor":        r.Flavor,
			"flavor_string": r.FlavorString,
			"method":        method,
			"url":           url,
			"reader":        reader,
		}).Debug("Sending HTTP request")
	}

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return result, err
	}

	defer func() {
		if request.Body != nil {
			ioutil.ReadAll(request.Body)
			request.Body.Close()
		}
	}()

	if basicAuth {
		if r.Debug {
			log.WithFields(log.Fields{
				"hostname":      r.Hostname,
				"port":          r.Port,
				"timeout":       r.Timeout,
				"flavor":        r.Flavor,
				"flavor_string": r.FlavorString,
				"method":        method,
				"url":           url,
			}).Debug("Setting HTTP basic authentication")
		}
		request.SetBasicAuth(r.Username, r.Password)
	}

	// add required headers
	request.Header.Add("OData-Version", "4.0") // Redfish API supports at least Open Data Protocol 4.0 (https://www.odata.org/documentation/)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	// set User-Agent
	request.Header.Set("User-Agent", userAgent)

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
		log.WithFields(log.Fields{
			"hostname":      r.Hostname,
			"port":          r.Port,
			"timeout":       r.Timeout,
			"flavor":        r.Flavor,
			"flavor_string": r.FlavorString,
			"method":        method,
			"url":           url,
			"http_headers":  request.Header,
		}).Debug("HTTP request headers")
	}

	response, err := client.Do(request)
	if err != nil {
		return result, err
	}

	defer func() {
		ioutil.ReadAll(response.Body)
		response.Body.Close()
	}()

	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":      r.Hostname,
			"port":          r.Port,
			"timeout":       r.Timeout,
			"flavor":        r.Flavor,
			"flavor_string": r.FlavorString,
			"method":        method,
			"url":           url,
			"status":        response.Status,
			"http_headers":  response.Header,
		}).Debug("HTTP reply received")
	}

	result.Status = response.Status
	result.StatusCode = response.StatusCode
	result.Header = response.Header
	result.Content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	if r.Debug {
		log.WithFields(log.Fields{
			"hostname":      r.Hostname,
			"port":          r.Port,
			"timeout":       r.Timeout,
			"flavor":        r.Flavor,
			"flavor_string": r.FlavorString,
			"method":        method,
			"url":           url,
			"content":       string(result.Content),
		}).Debug("Received content")
	}

	return result, nil
}
