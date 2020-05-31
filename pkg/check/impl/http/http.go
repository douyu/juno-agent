// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"encoding/json"
	"fmt"
	"github.com/armon/circbuf"
	"github.com/douyu/juno-agent/pkg/check/view"
	"time"

	"io"
	"net/http"
	"strings"
)

const (
	// DefaultBufSize is the maximum size of the captured
	// check output by default. Prevents an enormous buffer
	// from being captured
	DefaultBufSize = 4 * 1024
	// UserAgent is the value of the User-Agent header
	// for HTTP health checks.
	UserAgent = "Consul Health Check"
)

// HTTPHealthCheck http check config
type HTTPHealthCheck struct {
	httpClient *http.Client
	httpConfig
}

// httpConfig ...
type httpConfig struct {
	HTTP          string              `json:"http"`
	Header        map[string][]string `json:"header"`
	Method        string              `json:"method"`
	Body          string              `json:"body"`
	Interval      int64               `json:"interval"`
	Timeout       int64               `json:"timeout"`
	OutputMaxSize int                 `json:"output_max_size"`
	// Set if checks are exposed through Connect proxies
	// If set, this is the target of check()
	ProxyHTTP string `json:"proxy_http"`
}

// DefaultHTTPCheck return default http config
func DefaultHTTPHealthCheck() *HTTPHealthCheck {
	return &HTTPHealthCheck{
		httpConfig: httpConfig{
			Timeout:       3,
			OutputMaxSize: DefaultBufSize,
		},
		httpClient: &http.Client{
			Timeout: 3,
		},
	}
}

// NewHTTPHealthCheck new the instance
func NewHTTPHealthCheck() *HTTPHealthCheck {
	return DefaultHTTPHealthCheck()
}

// LoadExtConfig parse the config
func (c *HTTPHealthCheck) LoadExtConfig(extConfig string) (err error) {
	if err = json.Unmarshal([]byte(extConfig), &c.httpConfig); err != nil {
		return
	}

	c.httpClient.Timeout = time.Second * time.Duration(c.httpConfig.Timeout)
	return
}

// DoHealthCheck check is invoked periodically to perform the HTTP check
func (c *HTTPHealthCheck) DoHealthCheck() (resHealthCheck *view.ResHealthCheck, err error) {
	method := c.Method
	if method == "" {
		method = "GET"
	}
	target := c.HTTP
	if c.ProxyHTTP != "" {
		target = c.ProxyHTTP
	}

	bodyReader := strings.NewReader(c.Body)
	req, err := http.NewRequest(method, target, bodyReader)
	if err != nil {
		return resHealthCheck, err
	}

	req.Header = http.Header(c.Header)

	// this happens during testing but not in prod
	if req.Header == nil {
		req.Header = make(http.Header)
	}

	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", UserAgent)
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/plain, text/*, */*")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resHealthCheck, err
	}
	defer resp.Body.Close()

	// Read the response into a circular buffer to limit the size
	output, err := circbuf.NewBuffer(int64(c.OutputMaxSize))
	if _, err := io.Copy(output, resp.Body); err != nil {
		return resHealthCheck, err
	}

	// Format the response body
	result := fmt.Sprintf("HTTP %s %s: %s Output: %s", method, target, resp.Status, output.String())

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// PASSING (2xx)
		return view.HealthCheckResult("http", true, "success"), nil
	} else if resp.StatusCode == 429 {
		// WARNING
		// 429 Too Many Requests (RFC 6585)
		// The user has sent too many requests in a given amount of time.
		return view.HealthCheckResult("http", true, result), nil
	} else {
		// CRITICAL
		return view.HealthCheckResult("http", true, result), nil
	}
}
