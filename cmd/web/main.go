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

package main

import (
	"github.com/jmckind/grackle/pkg/grackle"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// listenAddress specifies the host:port on which to listen for requests.
	listenAddress = kingpin.Flag(
		"listen-address",
		"Host:port on which to listen for requests. Overrides GRK_LISTEN_ADDRESS.",
	).Default("0.0.0.0:8000").OverrideDefaultFromEnvar("GRK_LISTEN_ADDRESS").String()

	// logLevel specifies the logging level to use for the application.
	logLevel = kingpin.Flag(
		"log-level",
		"Logging level to use for Logrus. Overrides GRK_LOG_LEVEL.",
	).Default("info").OverrideDefaultFromEnvar("GRK_LOG_LEVEL").String()

	// rethinkDbHost specifies the RethinkDB hostname or IP address.
	rethinkDbHost = kingpin.Flag(
		"rethinkdb-host",
		"RethinkDB hostname or IP address. Overrides GRK_RETHINKDB_HOST.",
	).Default("localhost").OverrideDefaultFromEnvar("GRK_RETHINKDB_HOST").String()

	// rethinkDbPort specifies the RethinkDB port to use for the connection.
	rethinkDbPort = kingpin.Flag(
		"rethinkdb-port",
		"RethinkDB port for the connection. Overrides GRK_RETHINKDB_PORT.",
	).Default("28015").OverrideDefaultFromEnvar("GRK_RETHINKDB_PORT").String()

	// rethinkDbDatabase specifies the RethinkDB database name.
	rethinkDbDatabase = kingpin.Flag(
		"rethinkdb-database",
		"RethinkDB database name. Overrides GRK_RETHINKDB_DATABASE.",
	).Default("grackle").OverrideDefaultFromEnvar("GRK_RETHINKDB_DATABASE").String()

	// rethinkDbUsername specifies the RethinkDB username to use for the connection.
	rethinkDbUsername = kingpin.Flag(
		"rethinkdb-username",
		"RethinkDB username for the connection. Overrides GRK_RETHINKDB_USERNAME.",
	).OverrideDefaultFromEnvar("GRK_RETHINKDB_USERNAME").String()

	// rethinkDbPassword specifies the RethinkDB password to use for the connection.
	rethinkDbPassword = kingpin.Flag(
		"rethinkdb-password",
		"RethinkDB password for the connection. Overrides GRK_RETHINKDB_PASSWORD.",
	).OverrideDefaultFromEnvar("GRK_RETHINKDB_PASSWORD").String()

	// rethinkDbTLSCACertPath specifies the RethinkDB TLS CA certificate path.
	rethinkDbTLSCACertPath = kingpin.Flag(
		"rethinkdb-tls-ca-cert",
		"RethinkDB TLS CA certificate path. Overrides GRK_RETHINKDB_TLS_CA.",
	).OverrideDefaultFromEnvar("GRK_RETHINKDB_TLS_CA").String()

	// rethinkDbTLSClientCertPath specifies the RethinkDB TLS client certificate path.
	rethinkDbTLSClientCertPath = kingpin.Flag(
		"rethinkdb-tls-cert",
		"RethinkDB TLS client certificate path. Overrides GRK_RETHINKDB_TLS_CERT.",
	).OverrideDefaultFromEnvar("GRK_RETHINKDB_TLS_CERT").String()

	// rethinkDbTLSClientKeyPath specifies the RethinkDB client private key path.
	rethinkDbTLSClientKeyPath = kingpin.Flag(
		"rethinkdb-tls-key",
		"RethinkDB TLS client private key path. Overrides GRK_RETHINKDB_TLS_KEY.",
	).OverrideDefaultFromEnvar("GRK_RETHINKDB_TLS_KEY").String()

	templateDir = kingpin.Flag(
		"template-dir",
		"Directory for HTML templates used by the application. Overrides GRK_TEMPLATE_DIR.",
	).Default("templates").OverrideDefaultFromEnvar("GRK_TEMPLATE_DIR").String()
)

func main() {
	kingpin.Parse()
	opts := grackle.WebOptions{
		ListenAddress: *listenAddress,
		LogLevel:      *logLevel,
		RethinkDB: &grackle.RethinkdbOptions{
			Host:              *rethinkDbHost,
			Port:              *rethinkDbPort,
			Database:          *rethinkDbDatabase,
			Username:          *rethinkDbUsername,
			Password:          *rethinkDbPassword,
			TLSCACertPath:     *rethinkDbTLSCACertPath,
			TLSClientCertPath: *rethinkDbTLSClientCertPath,
			TLSClientKeyPath:  *rethinkDbTLSClientKeyPath,
		},
		TemplateDir: *templateDir,
	}

	app := grackle.NewWebApp(opts)
	app.Start()
}
