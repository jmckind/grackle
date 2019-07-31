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
	"fmt"
	"net/http"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

var (
	tweetsProcessed = promauto.NewCounter(prometheus.CounterOpts{
			Name: "grackle_processed_tweets_total",
			Help: "The total number of processed tweets",
	})
)

// IngestApp represents the Ingest process.
type IngestApp struct {
	log     *logrus.Logger
	opts    IngestOptions
	session *r.Session
	twitter *anaconda.TwitterApi
}

// IngestOptions stores data for tweet ingestion.
type IngestOptions struct {
	// LogLevel is the requested logging level for the application.
	LogLevel string

	// Twitter options.
	Twitter *TwitterOptions

	// RethinkDB options.
	RethinkDB *RethinkdbOptions
}

// NewIngestApp will create a new IngestApp.
func NewIngestApp(opts IngestOptions) *IngestApp {
	app := &IngestApp{}
	app.log = newLogger(opts.LogLevel)
	app.opts = opts
	app.opts.Twitter.Logger = newLogrusAnacondaLogger(app.log)
	return app
}

// Start will initialize and run the application.
func (a *IngestApp) Start() {
	a.log.Debug("starting application")
	logVersion(a.log)

	a.log.Debug("serving metrics")
	a.serveMetrics()

	a.log.Debug("connecting to rethinkdb")
	s, err := connectRethinkDB(a.opts.RethinkDB)
	if err != nil {
		a.log.Fatalf("unable to connect to rethibkdb: %v", err)
	}
	a.session = s

	a.log.Debug("connecting to twitter")
	t, err := connectTwitterAPI(a.opts.Twitter)
	if err != nil {
		a.log.Fatalf("unable to connect to twitter api: %v", err)
	}
	a.twitter = t

	a.log.Debug("ingesting stream")
	a.ingestStream()
}

// ingestStream will save a Tweet stream to a RethinkDB database.
func (a *IngestApp) ingestStream() {
	params := url.Values{
		"track": a.opts.Twitter.Track,
	}

	a.log.WithFields(logrus.Fields{"params": params}).Info("streaming tweets")
	stream := a.twitter.PublicStreamFilter(params)

	for obj := range stream.C {
		switch o := obj.(type) {
		case anaconda.Tweet:
			a.log.Tracef("%-15s: %s", o.User.ScreenName, o.Text)
			a.saveTweet(o)
			tweetsProcessed.Inc()
		}
	}
}

// saveTweet will save a Tweet in a RethinkDB database.
func (a *IngestApp) saveTweet(tweet anaconda.Tweet) {
	doc := make(map[string]anaconda.Tweet)
	doc["tweet"] = tweet

	err := r.Table("tweets").Insert(doc).Exec(a.session)
	if err != nil {
		a.log.Error("unable to save tweet", err)
	}
}

// serveMetrics will start a listening Prometheus metrics endpoint.
func (a *IngestApp) serveMetrics() {
	http.Handle(DefaultMetricsEdpoint, promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", DefaultMetricsPort), nil)
}
