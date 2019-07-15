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
	"strings"

	"github.com/jmckind/grackle/pkg/grackle"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
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

	// twitterAccessToken specifies the Twitter Access Token to use for the Twitter stream.
	twitterAccessToken = kingpin.Flag(
		"twitter-access-token",
		"Twitter Access Token  for the Twitter stream. Overrides GRK_TWITTER_ACCESS_TOKEN.",
	).OverrideDefaultFromEnvar("GRK_TWITTER_ACCESS_TOKEN").String()

	// twitterAccessSecret specifies the Twitter Access Secret to use for the Twitter stream.
	twitterAccessSecret = kingpin.Flag(
		"twitter-access-secret",
		"Twitter Access Secret for the Twitter stream. Overrides GRK_TWITTER_ACCESS_SECRET.",
	).OverrideDefaultFromEnvar("GRK_TWITTER_ACCESS_SECRET").String()

	// twitterConsumerKey specifies the Twitter Consumer Key to use for the Twitter stream.
	twitterConsumerKey = kingpin.Flag(
		"twitter-consumer-key",
		"Twitter Consumer Key for the Twitter stream. Overrides GRK_TWITTER_CONSUMER_KEY.",
	).OverrideDefaultFromEnvar("GRK_TWITTER_CONSUMER_KEY").String()

	// twitterConsumerSecret specifies the Twitter Consumer Secret to use for the Twitter stream.
	twitterConsumerSecret = kingpin.Flag(
		"twitter-consumer-secret",
		"Twitter Consumer Secret for the Twitter stream. Overrides GRK_TWITTER_CONSUMER_SECRET.",
	).OverrideDefaultFromEnvar("GRK_TWITTER_CONSUMER_SECRET").String()

	// twitterTrack specifies the search terms to use for the Twitter stream.
	twitterTrack = kingpin.Flag(
		"twitter-track",
		"Comma-delimited list of search terms for the Twitter stream. Overrides GRK_TWITTER_TRACK.",
	).OverrideDefaultFromEnvar("GRK_TWITTER_TRACK").Required().String()
)

func main() {
	kingpin.Parse()
	opts := grackle.IngestOptions{
		LogLevel: *logLevel,
		Twitter: &grackle.TwitterOptions{
			AccessToken:    *twitterAccessToken,
			AccessSecret:   *twitterAccessSecret,
			ConsumerKey:    *twitterConsumerKey,
			ConsumerSecret: *twitterConsumerSecret,
			Track:          strings.Split(*twitterTrack, ","),
		},
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
	}

	app := grackle.NewIngestApp(opts)
	app.Start()
}
