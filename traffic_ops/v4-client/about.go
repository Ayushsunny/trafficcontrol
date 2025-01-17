// Package client implements methods for interacting with the Traffic Ops API.
//
// Warning: Using the un-versioned import path ("client") is deprecated, and the
// ability to do so will be removed in ATC 6.0 - please use versioned client
// imports (e.g. "v3-client") instead
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
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APIAbout is the API version-relative path for the /about API endpoint.
	APIAbout = "/about"
)

// GetAbout gets data about the TO instance.
func (to *Session) GetAbout() (map[string]string, toclientlib.ReqInf, error) {
	route := APIAbout
	var data map[string]string
	reqInf, err := to.get(route, nil, &data)
	return data, reqInf, err
}
