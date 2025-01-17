package client

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

import (
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	apiCDNsDNSSECKeysGenerate = "/cdns/dnsseckeys/generate"
	apiCDNsNameDNSSECKeys     = "cdns/name/%s/dnsseckeys"
)

// GenerateCDNDNSSECKeys generates DNSSEC keys for the given CDN.
func (to *Session) GenerateCDNDNSSECKeys(req tc.CDNDNSSECGenerateReq, header http.Header) (tc.GenerateCDNDNSSECKeysResponse, toclientlib.ReqInf, error) {
	var resp tc.GenerateCDNDNSSECKeysResponse
	reqInf, err := to.post(apiCDNsDNSSECKeysGenerate, req, header, &resp)
	return resp, reqInf, err
}

// GetCDNDNSSECKeys gets the DNSSEC keys for the given CDN.
func (to *Session) GetCDNDNSSECKeys(name string, header http.Header) (tc.CDNDNSSECKeysResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf(apiCDNsNameDNSSECKeys, name)
	var resp tc.CDNDNSSECKeysResponse
	reqInf, err := to.get(route, header, &resp)
	return resp, reqInf, err
}

// DeleteCDNDNSSECKeys deletes all the DNSSEC keys for the given CDN.
func (to *Session) DeleteCDNDNSSECKeys(name string, header http.Header) (tc.DeleteCDNDNSSECKeysResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf(apiCDNsNameDNSSECKeys, name)
	var resp tc.DeleteCDNDNSSECKeysResponse
	reqInf, err := to.del(route, header, &resp)
	return resp, reqInf, err
}
