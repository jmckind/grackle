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
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	// EnvLogLevel is the log level environment variable.
	EnvLogLevel = "GRK_LOG_LEVEL"

	// DefaultMetricsEdpoint is the default endpoint name for application metrics.
	DefaultMetricsEdpoint = "/metrics"

	// DefaultMetricsPort is the default port for application metrics.
	DefaultMetricsPort = 8774
)

// newLogger will return a properly configured logger.
func newLogger(level string) *logrus.Logger {
	log := logrus.New()

	// Allow override of logging level.
	logLevel, err := logrus.ParseLevel(level)
	if err == nil {
		log.SetLevel(logLevel)
	}

	return log
}

// logVersion will log the current version.
func logVersion(log *logrus.Logger) {
	log.Debugf("Go Version: %s", runtime.Version())
	log.Debugf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Debugf("   Grackle: %s", Version)
}
