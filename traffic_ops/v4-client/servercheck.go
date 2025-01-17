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
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// APIServercheck is the API version-relative path to the /servercheck API endpoint.
const APIServercheck = "/servercheck"

// InsertServerCheckStatus will insert/update the Servercheck value based on if
// it already exists or not.
func (to *Session) InsertServerCheckStatus(status tc.ServercheckRequestNullable) (*tc.ServercheckPostResponse, toclientlib.ReqInf, error) {
	uri := APIServercheck
	resp := tc.ServercheckPostResponse{}
	reqInf, err := to.post(uri, status, nil, &resp)
	if err != nil {
		return nil, reqInf, err
	}
	return &resp, reqInf, nil
}

// GetServersChecks fetches check and meta information about servers from
// /servercheck.
func (to *Session) GetServersChecks(params url.Values, header http.Header) ([]tc.GenericServerCheck, tc.Alerts, toclientlib.ReqInf, error) {
	var response struct {
		tc.Alerts
		Response []tc.GenericServerCheck `json:"response"`
	}
	route := APIServercheck
	if params != nil {
		route += "?" + params.Encode()
	}
	reqInf, err := to.get(route, header, &response)
	return response.Response, response.Alerts, reqInf, err
}
