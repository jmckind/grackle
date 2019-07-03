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

// ConnectRethinkDB will return a Session object for the given RethinkDB database options.
func ConnectRethinkDB(opt *RethinkdbOptions) *r.Session {
	rdbOpts := r.ConnectOpts{
		Address:   fmt.Sprintf("%s:%s", opt.Host, opt.Port),
		Database:  opt.Database,
		Username:  opt.Username,
		Password:  opt.Password,
		TLSConfig: opt.TLSConfig,
	}

	session, err := r.Connect(rdbOpts)
	if err != nil {
		log.Fatalf("Unable to connect to database. %v", err)
	}

	err = r.DBCreate(opt.Database).Exec(session)
	if err != nil {
		log.Errorf("Unable to create database. %v", err)
	}

	err = r.TableCreate("tweets").Exec(session)
	if err != nil {
		log.Errorf("Unable to create table. %v", err)
	}
	return session
}
