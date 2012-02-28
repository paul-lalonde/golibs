// Copyright 2011 Twitter, Inc.
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

package twurlrc

import (
	"launchpad.net/goyaml"
	"os"
	"io/ioutil"
)

// Represents OAuth credentials to make requests on behalf of a user.
type Credentials struct {
	Token          string
	Username       string
	ConsumerKey    string
	ConsumerSecret string
	Secret         string
}

// Represents a parsed ~/.twurlrc formatted file.
type Twurlrc struct {
	data map[string]interface{}
}

// Loads a Twurlrc object from the standard ~/.twurlrc location.
func LoadTwurlrc() (*Twurlrc, os.Error) {
	t := new(Twurlrc)
	t.data = make(map[string]interface{})
	path := os.ShellExpand("$HOME/.twurlrc")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = goyaml.Unmarshal(data, t.data)
	return t, nil
}

// Returns credentials for the given user profile and consumer key.
func (t *Twurlrc) GetCredentials(profile string, key string) *Credentials {
	profileMap := t.data["profiles"].(map[interface{}]interface{})
	keyMap := profileMap[profile].(map[interface{}]interface{})
	data := keyMap[key].(map[interface{}]interface{})
	return &Credentials{
		Token:          data["token"].(string),
		Username:       data["username"].(string),
		ConsumerKey:    data["consumer_key"].(string),
		ConsumerSecret: data["consumer_secret"].(string),
		Secret:         data["secret"].(string),
	}
}

// Returns the default credentials, as specified in the ~/.twurlrc file.
func (t *Twurlrc) GetDefaultCredentials() *Credentials {
	configMap := t.data["configuration"].(map[interface{}]interface{})
	parts := configMap["default_profile"].([]interface{})
	return t.GetCredentials(parts[0].(string), parts[1].(string))
}

// Returns a list of consumer keys authorized with the given profile.
func (t *Twurlrc) GetKeys(profile string) []string {
	profileMap := t.data["profiles"].(map[string]interface{})
	keyMap := profileMap[profile].(map[string]interface{})
	keys := make([]string, len(keyMap))
	i := 0
	for key, _ := range keyMap {
		keys[i] = key
		i++
	}
	return keys
}

// Returns a list of profiles listed in the ~/.twurlrc file.
func (t *Twurlrc) GetProfiles() []string {
	profileMap := t.data["profiles"].(map[string]interface{})
	profiles := make([]string, len(profileMap))
	i := 0
	for key, _ := range profileMap {
		profiles[i] = key
		i++
	}
	return profiles
}
