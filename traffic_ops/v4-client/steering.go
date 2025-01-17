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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// Steering retrieves information about all (Tenant-accessible) Steering
// Delivery Services stored in Traffic Ops.
func (to *Session) Steering(header http.Header) ([]tc.Steering, toclientlib.ReqInf, error) {
	data := struct {
		Response []tc.Steering `json:"response"`
	}{}
	reqInf, err := to.get(`/steering`, header, &data)
	return data.Response, reqInf, err
}
