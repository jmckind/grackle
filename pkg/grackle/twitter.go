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
	"github.com/sirupsen/logrus"
)

// TwitterOptions stores data related to a Twitter connection.
type TwitterOptions struct {
	AccessToken    string
	AccessSecret   string
	ConsumerKey    string
	ConsumerSecret string
	Track          []string
	Logger         *logrusAnacondaLogger
}

// ConnectTwitterAPI will return a connection for the given Twitter options.
func connectTwitterAPI(opt *TwitterOptions) (*anaconda.TwitterApi, error) {
	api := anaconda.NewTwitterApiWithCredentials(
		opt.AccessToken,
		opt.AccessSecret,
		opt.ConsumerKey,
		opt.ConsumerSecret,
	)
	api.Log = *opt.Logger

	if ok, err := api.VerifyCredentials(); !ok || err != nil {
		return nil, err
	}
	return api, nil
}

type logrusAnacondaLogger struct {
	log *logrus.Logger
}

// NewLogrusAnacondaLogger constructs new logrusAnacondaLogger objects.
func newLogrusAnacondaLogger(log *logrus.Logger) *logrusAnacondaLogger {
	return &logrusAnacondaLogger{log: log}
}

func (l logrusAnacondaLogger) Fatal(items ...interface{})               { l.log.Fatal(items...) }
func (l logrusAnacondaLogger) Fatalf(s string, items ...interface{})    { l.log.Fatalf(s, items...) }
func (l logrusAnacondaLogger) Panic(items ...interface{})               { l.log.Panic(items...) }
func (l logrusAnacondaLogger) Panicf(s string, items ...interface{})    { l.log.Panicf(s, items...) }
func (l logrusAnacondaLogger) Critical(items ...interface{})            { l.log.Error(items...) }
func (l logrusAnacondaLogger) Criticalf(s string, items ...interface{}) { l.log.Errorf(s, items...) }
func (l logrusAnacondaLogger) Error(items ...interface{})               { l.log.Error(items...) }
func (l logrusAnacondaLogger) Errorf(s string, items ...interface{})    { l.log.Errorf(s, items...) }
func (l logrusAnacondaLogger) Warning(items ...interface{})             { l.log.Warn(items...) }
func (l logrusAnacondaLogger) Warningf(s string, items ...interface{})  { l.log.Warnf(s, items...) }
func (l logrusAnacondaLogger) Notice(items ...interface{})              { l.log.Info(items...) }
func (l logrusAnacondaLogger) Noticef(s string, items ...interface{})   { l.log.Infof(s, items...) }
func (l logrusAnacondaLogger) Info(items ...interface{})                { l.log.Info(items...) }
func (l logrusAnacondaLogger) Infof(s string, items ...interface{})     { l.log.Infof(s, items...) }
func (l logrusAnacondaLogger) Debug(items ...interface{})               { l.log.Debug(items...) }
func (l logrusAnacondaLogger) Debugf(s string, items ...interface{})    { l.log.Debugf(s, items...) }
