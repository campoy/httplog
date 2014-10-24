// Copyright 2014 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

// httplog provides an implementation of http.RoundTripper that logs every
// single request and response using a given logging function.
package httplog

import (
	"log"
	"net/http"
	"net/http/httputil"
)

// Transport logs all requests and responses at every RoundTrip using the
// provided logging function.
type Transport struct {
	Transport http.RoundTripper
	// Should the body of the requests and responses be logged.
	LogBody bool
	// If Logf is nil log.Printf will be used.
	Logf func(format string, vs ...interface{})
}

// Transport satifies http.RoundTripper
func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Logf == nil {
		t.Logf = log.Printf
	}
	if b, err := httputil.DumpRequest(req, t.LogBody); err != nil {
		t.Logf("httplog: dump request: %v", err)
		return nil, err
	} else {
		t.Logf("httplog: %s", b)
	}

	res, err := t.Transport.RoundTrip(req)
	if err != nil {
		t.Logf("httplog: roundtrip error: %v", err)
		return res, err
	}

	if b, err := httputil.DumpResponse(res, t.LogBody); err != nil {
		t.Logf("httplog: dump response: %v", err)
	} else {
		t.Logf("httplog: %s", b)
	}
	return res, err
}
