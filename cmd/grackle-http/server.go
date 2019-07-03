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
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jmckind/grackle"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var tmpl *template.Template
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ServerOpts stores configuration options for the server.
type ServerOpts struct {
	// Addr is the host and port to listen for requests.
	Addr string

	// TemplateDir is the directory that contains HTML templates.
	TemplateDir string
}

func main() {
	log.SetLevel(logrus.DebugLevel)
	grackle.PrintVersion()

	startServer(ServerOpts{
		Addr:        "0.0.0.0:8000",
		TemplateDir: "templates",
	})
}

func startServer(opts ServerOpts) {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/sock", socketHandler)

	tmpl = template.Must(template.ParseGlob(fmt.Sprintf("%s/index.html.tmpl", opts.TemplateDir)))
	srv := &http.Server{
		Handler:      r,
		Addr:         opts.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("Listening at %s", opts.Addr)
	log.Fatal(srv.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{"address": r.RemoteAddr}).Debugf("Request for index")
	tmpl.Execute(w, nil)
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{"address": r.RemoteAddr}).Debugf("Request for sock")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Connection upgraded...")
	// Use connection...
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}

		if messageType != websocket.TextMessage {
			log.Warnf("Non-text message received: %v", messageType)
			return
		}

		log.Debugf("Receieved message: %v", p)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Error(err)
			return
		}
	}
}
