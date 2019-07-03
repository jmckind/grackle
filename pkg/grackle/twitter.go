// Copyright 2019 Grackle Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grackle

import (
	"github.com/ChimeraCoder/anaconda"
)

var log = anaconda.BasicLogger

// TwitterOptions stores data related to a Twitter connection.
type TwitterOptions struct {
	AccessToken    string
	AccessSecret   string
	ConsumerKey    string
	ConsumerSecret string
	Track          []string
}

// ConnectTwitterAPI will return a connection for the given Twitter options.
func ConnectTwitterAPI(opt *TwitterOptions) *anaconda.TwitterApi {
	api := anaconda.NewTwitterApiWithCredentials(
		opt.AccessToken,
		opt.AccessSecret,
		opt.ConsumerKey,
		opt.ConsumerSecret,
	)
	api.Log = log

	if ok, err := api.VerifyCredentials(); !ok || err != nil {
		log.Fatalf("Invalid credentials. %v", err)
	}
	return api
}
