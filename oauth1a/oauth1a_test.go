// Copyright 2011 Arne Roomann-Kurrik.
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

package oauth1a

import (
	"testing"
	"http"
)

var user = NewAuthorizedConfig("token", "secret")

var client = &ClientConfig{
	ConsumerKey:    "consumer_key",
	ConsumerSecret: "consumer_secret",
	CallbackURL:    "https://example.com/callback",
}

var signer = new(HmacSha1Signer)

var service = &Service{
	RequestURL:   "https://example.com/request_token",
	AuthorizeURL: "https://example.com/request_token",
	AccessURL:    "https://example.com/request_token",
	ClientConfig: client,
	Signer:       signer,
}

func TestSignature(t *testing.T) {
	url := "https://example.com/endpoint"
	request, _ := http.NewRequest("GET", url, nil)
	service.Sign(request, user)
	params, _ := signer.GetOAuthParams(request, client, user, "nonce", "timestamp")
	signature := params["oauth_signature"]
	expected := "8+ZC6DP8FU3z50qSWDeYCGix2x0="
	if signature != expected {
		t.Errorf("Signature %v did not match expected %v", signature, expected)
	}
}

var ESCAPE_TESTS = map[string]string{
	"Ā" : "%C4%80",
        "㤹" : "%E3%A4%B9",
}

func TestEscaping(t *testing.T) {
	for str, expected := range(ESCAPE_TESTS) {
		if Rfc3986Escape(str) != expected {
			t.Errorf("Escaped %v was %v, expected %v", str, Rfc3986Escape(str), expected)
		}
	}
}
