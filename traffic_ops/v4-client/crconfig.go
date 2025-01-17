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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APISnaphost is the API version-relative path for the /snapshot API endpoint.
	APISnaphost = "/snapshot"
)

// OuterResponse is the most basic type of a Traffic Ops API response,
// containing some kind of JSON-encoded 'response' object.
type OuterResponse struct {
	Response json.RawMessage `json:"response"`
}

// GetCRConfig returns the raw JSON bytes of the CRConfig from Traffic Ops, and whether the bytes were from the client's internal cache.
func (to *Session) GetCRConfig(cdn string) ([]byte, toclientlib.ReqInf, error) {
	uri := `/cdns/` + cdn + `/snapshot`
	bts := []byte{}
	reqInf, err := to.get(uri, nil, &bts)
	if err != nil {
		return nil, reqInf, err
	}

	resp := OuterResponse{}
	if err := json.Unmarshal(bts, &resp); err != nil {
		return nil, reqInf, err
	}
	return resp.Response, reqInf, nil
}

// SnapshotCRConfig creates a new Snapshot for the CDN with the given Name -
// NOT just a new CRConfig!
func (to *Session) SnapshotCRConfig(cdn string, header http.Header) (toclientlib.ReqInf, error) {
	uri := fmt.Sprintf("%s?cdn=%s", APISnaphost, url.QueryEscape(cdn))
	resp := OuterResponse{}
	reqInf, err := to.put(uri, nil, header, &resp)
	return reqInf, err
}

// GetCRConfigNew returns the raw JSON bytes of the latest CRConfig from Traffic Ops, and whether the bytes were from the client's internal cache.
func (to *Session) GetCRConfigNew(cdn string) ([]byte, toclientlib.ReqInf, error) {
	uri := `/cdns/` + cdn + `/snapshot/new`
	bts := []byte{}
	reqInf, err := to.get(uri, nil, &bts)
	if err != nil {
		return nil, reqInf, err
	}

	resp := OuterResponse{}
	if err := json.Unmarshal(bts, &resp); err != nil {
		return nil, reqInf, err
	}
	return resp.Response, reqInf, nil
}

// SnapshotCRConfigByID creates a new Snapshot for the CDN with the given ID -
// NOT just a new CRConfig!
func (to *Session) SnapshotCRConfigByID(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	url := fmt.Sprintf("%s?cdnID=%d", APISnaphost, id)
	var alerts tc.Alerts
	reqInf, err := to.put(url, nil, nil, &alerts)
	return alerts, reqInf, err
}
