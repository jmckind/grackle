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
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// WebApp represents the web ui process.
type WebApp struct {
	log      *logrus.Logger
	opts     WebOptions
	session  *r.Session
	tmpl     *template.Template
	upgrader websocket.Upgrader
}

// WebOptions stores configuration options for the server.
type WebOptions struct {
	// LogLevel is the requested logging level for the application.
	LogLevel string

	// ListenAddress is the host:port to listen for requests.
	ListenAddress string

	// TemplateDir is the directory that contains HTML templates.
	TemplateDir string

	// RethinkDB options.
	RethinkDB *RethinkdbOptions
}

// NewWebApp will create a new WebApp.
func NewWebApp(opts WebOptions) *WebApp {
	app := &WebApp{}
	app.log = newLogger(opts.LogLevel)
	app.opts = opts
	return app
}

// Start will initialize and run the application.
func (a *WebApp) Start() {
	a.log.Debug("starting application")
	logVersion(a.log)

	a.log.Debug("parsing index template")
	a.tmpl = template.Must(template.ParseGlob(fmt.Sprintf("%s/index.html.tmpl", a.opts.TemplateDir)))

	a.log.Debug("connecting to rethinkdb")
	s, err := connectRethinkDB(a.opts.RethinkDB)
	if err != nil {
		a.log.Fatalf("unable to connect to rethibkdb: %v", err)
	}
	a.session = s

	a.log.Debug("starting web server")
	a.startServer()
}

func (a *WebApp) startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", a.indexHandler)
	r.HandleFunc("/sock", a.socketHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         a.opts.ListenAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	a.log.Infof("listening at %s", a.opts.ListenAddress)
	a.log.Fatal(srv.ListenAndServe())
}

func (a *WebApp) indexHandler(w http.ResponseWriter, r *http.Request) {
	a.log.WithFields(logrus.Fields{"address": r.RemoteAddr}).Debugf("Request for index")
	a.tmpl.Execute(w, nil)
}

func (a *WebApp) socketHandler(w http.ResponseWriter, r *http.Request) {
	a.log.WithFields(logrus.Fields{"address": r.RemoteAddr}).Debugf("upgrading connection for websockets")
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.log.Error(err)
		return
	}
	defer conn.Close()

	// Use websocket connection...
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			a.log.Error(err)
			return
		}
		if messageType != websocket.TextMessage {
			a.log.Warnf("Non-text message received: %v", messageType)
			return
		}
		a.log.Debugf("Receieved message: %s", p)

		a.streamTweets(conn)
	}
}

func (a *WebApp) streamTweets(conn *websocket.Conn) {
	res, err := r.Table("tweets").Changes().Run(a.session)
	if err != nil {
		a.log.Error("unable to subscribe to changefeed", err)
		return
	}

	var change map[string]interface{}
	for res.Next(&change) {
		if change["old_val"] != nil || change["new_val"] == nil {
			log.Warn("skipping existing tweet...")
			continue
		}

		// Write new tweet to websocket...
		tweet, err := json.Marshal(change["new_val"])
		if err != nil {
			a.log.Error("unable to marshal change", err)
		}
		a.log.Tracef("new tweet: %s", tweet)

		if err = conn.WriteMessage(websocket.TextMessage, tweet); err != nil {
			a.log.Error(err)
		}
	}
	if res.Err() != nil {
		a.log.Error("unable to parse changefeed response", err)
	}
}

// reverse accepts a slice of bytes and returns the sequence reversed.
func reverse(b []byte) (r []byte) {
	n := len(b)
	r = make([]byte, n)
	for i := 0; i < n; i++ {
		r[i] = b[(n-i)-1]
	}
	return
}
