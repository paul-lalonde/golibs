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
	"io/ioutil"
	"http"
	"os"
)

// Container for user-specific keys and secrets related to the OAuth process.
// This struct is intended to be serialized and stored for future use.
// Request and Access tokens are each stored separately, so that the current
// position in the auth flow may be inferred.
type UserConfig struct {
	RequestTokenSecret string
	RequestTokenKey    string
	AccessTokenSecret  string
	AccessTokenKey     string
	Verifier           string
	AccessValues       http.Values
}

// Creates a UserConfig object with existing access token credentials.  For
// users where an access token has been obtained through other means than
// the authz flows provided by this library.
func NewAuthorizedConfig(token string, secret string) *UserConfig {
	return &UserConfig{AccessTokenKey: token, AccessTokenSecret: secret}
}

// Sign and send a Request using the current configuration.
func (c *UserConfig) send(request *http.Request, service *Service, client *http.Client) (*http.Response, os.Error) {
	if err := service.Sign(request, c); err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, os.NewError("Endpoint response: " + response.Status)
	}
	return response, nil
}

// Issue a request to obtain a Request token.
func (c *UserConfig) GetRequestToken(service *Service, client *http.Client) os.Error {
	request, err := http.NewRequest("POST", service.RequestURL, nil)
	if err != nil {
		return err
	}
	request.Form.Add("oauth_callback", service.ClientConfig.CallbackURL)
	response, err := c.send(request, service, client)
	if err != nil {
		return err
	}
	err = c.parseRequestToken(response)
	return err
}

// Given the returned response from a Request token request, parse out the
// appropriate request token and secret fields.
func (c *UserConfig) parseRequestToken(response *http.Response) os.Error {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	params, err := http.ParseQuery(string(body))
	tokenKey := params.Get("oauth_token")
	tokenSecret := params.Get("oauth_token_secret")
	if tokenKey == "" || tokenSecret == "" {
		return os.NewError("No token or secret found")
	}
	c.RequestTokenKey = tokenKey
	c.RequestTokenSecret = tokenSecret
	if params.Get("oauth_callback_confirmed") == "false" {
		return os.NewError("OAuth callback not confirmed")
	}
	return nil
}

// Obtain a URL which will allow the current user to authorize access to their
// OAuth-protected data.
func (c *UserConfig) GetAuthorizeURL(service *Service) (string, os.Error) {
	if c.RequestTokenKey == "" || c.RequestTokenSecret == "" {
		return "", os.NewError("No configured request token")
	}
	token := http.URLEscape(c.RequestTokenKey)
	return service.AuthorizeURL + "?oauth_token=" + token, nil
}

// Issue a request to exchange the current request token for an access token.
func (c *UserConfig) GetAccessToken(token string, verifier string, service *Service, client *http.Client) os.Error {
	if c.RequestTokenKey == "" || c.RequestTokenSecret == "" {
		return os.NewError("No configured request token")
	}
	if c.RequestTokenKey != token {
		return os.NewError("Returned token did not match request token")
	}
	c.Verifier = verifier
	request, err := http.NewRequest("POST", service.AccessURL, nil)
	if err != nil {
		return err
	}
	request.Form.Add("oauth_verifier", verifier)
	response, err := c.send(request, service, client)
	if err != nil {
		return err
	}
	err = c.parseAccessToken(response)
	return err
}

// Given the returned response from the access token request, pull out the
// access token and token secret.  Store a copy of any other values returned,
// too, since some services (like Twitter) return handy information such
// as the username.
func (c *UserConfig) parseAccessToken(response *http.Response) os.Error {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	params, err := http.ParseQuery(string(body))
	tokenKey := params.Get("oauth_token")
	tokenSecret := params.Get("oauth_token_secret")
	if tokenKey == "" || tokenSecret == "" {
		return os.NewError("No token or secret found")
	}
	c.AccessTokenKey = tokenKey
	c.AccessTokenSecret = tokenSecret
	c.AccessValues = params
	return nil
}

// Returns a token and secret corresponding to where in the OAuth flow this
// config is currently in.  The priority is Access token, Request token, empty
// string.
func (c *UserConfig) GetToken() (string, string) {
	if c.AccessTokenKey != "" {
		return c.AccessTokenKey, c.AccessTokenSecret
	}
	if c.RequestTokenKey != "" {
		return c.RequestTokenKey, c.RequestTokenSecret
	}
	return "", ""
}


