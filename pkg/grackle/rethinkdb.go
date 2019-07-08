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
	"crypto/tls"
	"fmt"

	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// RethinkdbOptions stores data related to a RethinkDB connection.
type RethinkdbOptions struct {
	Host              string
	Port              string
	Database          string
	Username          string
	Password          string
	TLSCACertPath     string
	TLSClientCertPath string
	TLSClientKeyPath  string
	TLSConfig         *tls.Config
}

// connectRethinkDB will return a Session object for the given RethinkDB database options.
func connectRethinkDB(opt *RethinkdbOptions) (*r.Session, error) {
	rdbOpts := r.ConnectOpts{
		Address:   fmt.Sprintf("%s:%s", opt.Host, opt.Port),
		Database:  opt.Database,
		Username:  opt.Username,
		Password:  opt.Password,
		TLSConfig: opt.TLSConfig,
	}

	session, err := r.Connect(rdbOpts)
	if err != nil {
		return nil, err
	}

	// Ensure database and table are created.
	_ = r.DBCreate(opt.Database).Exec(session)
	_ = r.TableCreate("tweets").Exec(session)

	return session, nil
}
