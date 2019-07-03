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
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/jmckind/grackle"
	"github.com/sirupsen/logrus"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

var log = logrus.New()

// EnvOptions stores data from the envioronment.
type EnvOptions struct {
	Twitter   *grackle.TwitterOptions
	RethinkDB *grackle.RethinkdbOptions
}

func main() {
	grackle.PrintVersion()
	archiveTweets()
}

func archiveTweets() {
	env := processEnvironment()
	twitter := grackle.connectTwitterAPI(env.Twitter)
	session := grackle.ConnectRethinkDB(env.RethinkDB)

	params := url.Values{
		"track": env.Twitter.Track,
	}

	log.Debugf("Streaming tweets using params: %v", params)
	stream := twitter.PublicStreamFilter(params)

	for obj := range stream.C {
		switch o := obj.(type) {
		case anaconda.Tweet:
			log.Debugf("%-15s: %s", o.User.ScreenName, o.Text)
			err := r.Table("tweets").Insert(o).Exec(session)
			if err != nil {
				log.Errorf("Unable to insert database record. %v", err)
			}
		}
	}
}

func processEnvironment() *EnvOptions {
	return &EnvOptions{
		Twitter:   processTwitterOptions(),
		RethinkDB: processRethinkDBOptions(),
	}
}

func processRethinkDBOptions() *RethinkdbOptions {
	opt := &grackle.RethinkdbOptions{
		Host:              os.Getenv("GTR_RETHINKDB_HOST"),
		Port:              os.Getenv("GTR_RETHINKDB_PORT"),
		Database:          os.Getenv("GTR_RETHINKDB_DATABASE"),
		Username:          os.Getenv("GTR_RETHINKDB_USERNAME"),
		Password:          os.Getenv("GTR_RETHINKDB_PASSWORD"),
		TLSCACertPath:     os.Getenv("GTR_RETHINKDB_TLS_CA"),
		TLSClientCertPath: os.Getenv("GTR_RETHINKDB_TLS_CERT"),
		TLSClientKeyPath:  os.Getenv("GTR_RETHINKDB_TLS_KEY"),
	}

	if len(opt.TLSCACertPath) > 0 && len(opt.TLSClientCertPath) > 0 {
		certPool := x509.NewCertPool()
		caCert, err := ioutil.ReadFile(opt.TLSCACertPath)
		if err != nil {
			log.Fatalf("Unable to parse CA certificate. %v", err)
		}
		certPool.AppendCertsFromPEM(caCert)

		clientCert, err := tls.LoadX509KeyPair(opt.TLSClientCertPath, opt.TLSClientKeyPath)
		if err != nil {
			log.Fatalf("Unable to parse client key pair. %v", err)
		}

		opt.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{clientCert},
			RootCAs:      certPool,
		}
	}

	log.Debugf("RethinkDB Host: %s", opt.Host)
	log.Debugf("RethinkDB Port: %s", opt.Port)
	log.Debugf("RethinkDB Database: %s", opt.Database)
	log.Debugf("RethinkDB Username: %s", opt.Username)
	log.Debugf("RethinkDB Password: %s", opt.Password)
	log.Debugf("RethinkDB TLS CA Path: %s", opt.TLSCACertPath)
	log.Debugf("RethinkDB TLS Cert Path: %s", opt.TLSClientCertPath)
	log.Debugf("RethinkDB TLS Key Path: %s", opt.TLSClientKeyPath)

	return opt
}

func processTwitterOptions() *TwitterOptions {
	return &TwitterOptions{
		AccessToken:    os.Getenv("GTR_TWITTER_ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("GTR_TWITTER_ACCESS_SECRET"),
		ConsumerKey:    os.Getenv("GTR_TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("GTR_TWITTER_CONSUMER_SECRET"),
		Track:          []string{os.Getenv("GTR_TWITTER_TRACK")},
	}
}
