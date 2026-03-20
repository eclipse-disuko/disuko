// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package requestHelper

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/eclipse-disuko/disuko/logy"
)

func DoPostFormRequestWithBasicAuth(endpoint string, header map[string]string, data url.Values, username string, password string, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode())) // URL-encoded payload

	if err != nil {
		return nil, err
	}

	if len(username) > 0 {
		req.SetBasicAuth(url.QueryEscape(username), url.QueryEscape(password))
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func DoPostFormRequest(endpoint string, header map[string]string, data url.Values, client *http.Client) ([]byte, error) {
	return DoPostFormRequestWithBasicAuth(endpoint, header, data, "", "", client)
}

func DoPostJSONRequest(requestSession *logy.RequestSession, endpoint string, jsonStr string, header map[string]string, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(jsonStr))
	if err != nil {
		logy.Errorf(requestSession, "Error creating request for endpoint %s", endpoint)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(jsonStr)))
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		logy.Errorf(requestSession, "Error requesting endpoint %s", endpoint)
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		logy.Errorf(requestSession, "Error reading answer for endpoint %s", endpoint)
		return nil, err
	}

	if res.StatusCode >= 400 {
		logy.Errorf(requestSession, "Error requesting endpoint %s. Response status %s", endpoint, res.Status)
		return nil, errors.New("Error requesting endpoint " + endpoint + ". Response status " + res.Status)
	}

	return response, nil
}

func DoGetRequest(requestSession *logy.RequestSession, endpoint string, header map[string]string, queryData url.Values, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		logy.Errorf(requestSession, "Error creating GET request for endpoint %s", endpoint)
		return nil, err
	}

	if queryData != nil {
		query := req.URL.Query()
		for key, values := range queryData {
			query.Set(key, values[0])
		}
		req.URL.RawQuery = query.Encode()
	}

	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		logy.Errorf(requestSession, "Error requesting endpoint %s", endpoint)
		return nil, err
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		logy.Errorf(requestSession, "Error reading answer for endpoint %s", endpoint)
		return nil, err
	}

	if res.StatusCode >= 400 {
		logy.Errorf(requestSession, "Error requesting endpoint %s. Response status %s", endpoint, res.Status)
		return nil, errors.New("Error requesting endpoint " + endpoint + ". Response status " + res.Status)
	}

	return response, nil
}
