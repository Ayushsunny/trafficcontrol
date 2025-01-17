/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APIServerCapabilities is the full path to the /server_capabilities API
	// endpoint.
	APIServerCapabilities = "/server_capabilities"
)

// CreateServerCapability creates the given Server Capability.
func (to *Session) CreateServerCapability(sc tc.ServerCapability) (*tc.ServerCapabilityDetailResponse, toclientlib.ReqInf, error) {
	var scResp tc.ServerCapabilityDetailResponse
	reqInf, err := to.post(APIServerCapabilities, sc, nil, &scResp)
	if err != nil {
		return nil, reqInf, err
	}
	return &scResp, reqInf, nil
}

// GetServerCapabilities returns all the Server Capabilities in Traffic Ops.
func (to *Session) GetServerCapabilities(header http.Header) ([]tc.ServerCapability, toclientlib.ReqInf, error) {
	var data tc.ServerCapabilitiesResponse
	reqInf, err := to.get(APIServerCapabilities, header, &data)
	return data.Response, reqInf, err
}

// GetServerCapability retrieves the Server Capability with the given Name.
func (to *Session) GetServerCapability(name string, header http.Header) (*tc.ServerCapability, toclientlib.ReqInf, error) {
	reqURL := fmt.Sprintf("%s?name=%s", APIServerCapabilities, url.QueryEscape(name))
	var data tc.ServerCapabilitiesResponse
	reqInf, err := to.get(reqURL, header, &data)
	if err != nil {
		return nil, reqInf, err
	}
	if len(data.Response) == 1 {
		return &data.Response[0], reqInf, nil
	}
	return nil, reqInf, fmt.Errorf("expected one server capability in response, instead got: %+v", data.Response)
}

// UpdateServerCapability updates a Server Capability by name.
func (to *Session) UpdateServerCapability(name string, sc *tc.ServerCapability, header http.Header) (*tc.ServerCapability, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", APIServerCapabilities, url.QueryEscape(name))
	var data tc.ServerCapability
	reqInf, err := to.put(route, sc, header, &data)
	if err != nil {
		return nil, reqInf, err
	}
	return &data, reqInf, nil
}

// DeleteServerCapability deletes the given server capability by name.
func (to *Session) DeleteServerCapability(name string) (tc.Alerts, toclientlib.ReqInf, error) {
	reqURL := fmt.Sprintf("%s?name=%s", APIServerCapabilities, url.QueryEscape(name))
	var alerts tc.Alerts
	reqInf, err := to.del(reqURL, nil, &alerts)
	return alerts, reqInf, err
}
