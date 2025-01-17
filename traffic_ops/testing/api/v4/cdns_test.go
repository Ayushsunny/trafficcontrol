package v4

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
	"net/http"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestCDNs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Parameters}, func() {
		GetTestCDNsIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		SortTestCDNs(t)
		UpdateTestCDNs(t)
		UpdateTestCDNsWithHeaders(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestCDNsWithHeaders(t, header)
		GetTestCDNs(t)
		GetTestCDNsIMSAfterChange(t, header)
	})
}

func TestCDNsDNSSEC(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServerCapabilities, DeliveryServices}, func() {
		if includeSystemTests {
			GenerateDNSSECKeys(t)
		}
	})
}

func GenerateDNSSECKeys(t *testing.T) {
	if len(testData.CDNs) < 1 {
		t.Fatalf("need at least one CDN to test updating CDNs")
	}
	firstCDN := testData.CDNs[0]
	resp, _, err := TOSession.GetCDNByName(firstCDN.Name, nil)
	if err != nil {
		t.Fatalf("cannot GET CDN by name: '%s', %v", firstCDN.Name, err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected: 1 CDN, actual: %d", len(resp))
	}

	ttl := util.JSONIntStr(60)
	req := tc.CDNDNSSECGenerateReq{
		Key:               util.StrPtr(firstCDN.Name),
		TTL:               &ttl,
		KSKExpirationDays: &ttl,
		ZSKExpirationDays: &ttl,
	}
	_, _, err = TOSession.GenerateCDNDNSSECKeys(req, nil)
	if err != nil {
		t.Fatalf("generating CDN DNSSEC keys - expected: nil error, actual: %s", err.Error())
	}

	res, _, err := TOSession.GetCDNDNSSECKeys(firstCDN.Name, nil)
	if err != nil {
		t.Fatalf("getting CDN DNSSEC keys - expected: nil error, actual: %s", err.Error())
	}
	if _, ok := res.Response[firstCDN.Name]; !ok {
		t.Errorf("getting CDN DNSSEC keys - expected: key %s, actual: missing", firstCDN.Name)
	}
	originalKeys := res.Response

	_, _, err = TOSession.GenerateCDNDNSSECKeys(req, nil)
	if err != nil {
		t.Fatalf("generating CDN DNSSEC keys - expected: nil error, actual: %s", err.Error())
	}
	res, _, err = TOSession.GetCDNDNSSECKeys(firstCDN.Name, nil)
	if err != nil {
		t.Fatalf("getting CDN DNSSEC keys - expected: nil error, actual: %s", err.Error())
	}
	newKeys := res.Response

	if reflect.DeepEqual(originalKeys, newKeys) {
		t.Errorf("generating CDN DNSSEC keys - expected: original keys to differ from new keys, actual: they are the same")
	}

	_, _, err = TOSession.DeleteCDNDNSSECKeys(firstCDN.Name, nil)
	if err != nil {
		t.Errorf("deleting CDN DNSSEC keys - expected: nil error, actual: %s", err.Error())
	}
}

func UpdateTestCDNsWithHeaders(t *testing.T, header http.Header) {
	if len(testData.CDNs) < 1 {
		t.Fatalf("need at least one CDN to test updating CDNs")
	}
	firstCDN := testData.CDNs[0]
	// Retrieve the CDN by name so we can get the id for the Update
	resp, _, err := TOSession.GetCDNByName(firstCDN.Name, header)
	if err != nil {
		t.Errorf("cannot GET CDN by name: '%s', %v", firstCDN.Name, err)
	}
	if len(resp) > 0 {
		remoteCDN := resp[0]
		remoteCDN.DomainName = "domain2"
		_, reqInf, err := TOSession.UpdateCDN(remoteCDN.ID, remoteCDN, header)
		if err == nil {
			t.Errorf("Expected error about Precondition Failed, got none")
		}
		if reqInf.StatusCode != http.StatusPreconditionFailed {
			t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestCDNsIMSAfterChange(t *testing.T, header http.Header) {
	for _, cdn := range testData.CDNs {
		_, reqInf, err := TOSession.GetCDNByName(cdn.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, cdn := range testData.CDNs {
		_, reqInf, err := TOSession.GetCDNByName(cdn.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestCDNsIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	for _, cdn := range testData.CDNs {
		futureTime := time.Now().AddDate(0, 0, 1)
		time := futureTime.Format(time.RFC1123)
		header.Set(rfc.IfModifiedSince, time)
		_, reqInf, err := TOSession.GetCDNByName(cdn.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.CreateCDN(cdn)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE cdns: %v", err)
		}
	}

}

func SortTestCDNs(t *testing.T) {
	var header http.Header
	var sortedList []string
	resp, _, err := TOSession.GetCDNs(header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	for i := range resp {
		sortedList = append(sortedList, resp[i].Name)
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func UpdateTestCDNs(t *testing.T) {

	firstCDN := testData.CDNs[0]
	// Retrieve the CDN by name so we can get the id for the Update
	resp, _, err := TOSession.GetCDNByName(firstCDN.Name, nil)
	if err != nil {
		t.Errorf("cannot GET CDN by name: '%s', %v", firstCDN.Name, err)
	}
	remoteCDN := resp[0]
	expectedCDNDomain := "domain2"
	remoteCDN.DomainName = expectedCDNDomain
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateCDN(remoteCDN.ID, remoteCDN, nil)
	if err != nil {
		t.Errorf("cannot UPDATE CDN by id: %v - %v", err, alert)
	}

	// Retrieve the CDN to check CDN name got updated
	resp, _, err = TOSession.GetCDNByID(remoteCDN.ID, nil)
	if err != nil {
		t.Errorf("cannot GET CDN by name: '$%s', %v", firstCDN.Name, err)
	}
	respCDN := resp[0]
	if respCDN.DomainName != expectedCDNDomain {
		t.Errorf("results do not match actual: %s, expected: %s", respCDN.DomainName, expectedCDNDomain)
	}

}

func GetTestCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.GetCDNByName(cdn.Name, nil)
		if err != nil {
			t.Errorf("cannot GET CDN by name: %v - %v", err, resp)
		}
	}
}

func DeleteTestCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		// Retrieve the CDN by name so we can get the id for the Update
		resp, _, err := TOSession.GetCDNByName(cdn.Name, nil)
		if err != nil {
			t.Errorf("cannot GET CDN by name: %v - %v", cdn.Name, err)
		}
		if len(resp) > 0 {
			respCDN := resp[0]

			_, _, err := TOSession.DeleteCDN(respCDN.ID)
			if err != nil {
				t.Errorf("cannot DELETE CDN by name: '%s' %v", respCDN.Name, err)
			}

			// Retrieve the CDN to see if it got deleted
			cdns, _, err := TOSession.GetCDNByName(cdn.Name, nil)
			if err != nil {
				t.Errorf("error deleting CDN name: %s", err.Error())
			}
			if len(cdns) > 0 {
				t.Errorf("expected CDN name: %s to be deleted", cdn.Name)
			}
		}
	}
}
