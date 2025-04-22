// Copyright 2023 Northern.tech AS
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
package log

import (
	"fmt"
	"log/slog"
	"os"
)

var (
	logOpts = slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger = slog.New(
		slog.NewTextHandler(os.Stderr, &logOpts),
	)
)

func Setup(verb bool) {
	if verb {
		logOpts.Level = slog.LevelDebug
		logger = slog.New(slog.NewTextHandler(os.Stderr, &logOpts))
	}
}

func Err(msg string) {
	logger.Error(msg)
}

func Errf(msg string, args ...interface{}) {
	logger.Error(fmt.Sprintf(msg, args...))
}

func Verb(msg string) {
	logger.Debug(msg)
}

func Verbf(msg string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(msg, args...))
}

func Info(msg string) {
	logger.Info(msg)
}

func Infof(msg string, args ...interface{}) {
	logger.Info(fmt.Sprintf(msg, args...))
}
