//Package gologging provides a logger implementation based on the github.com/op/go-logging pkg
package gologging

import (
	"fmt"
	"io"

	"github.com/MagalixTechnologies/core/logger"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
)

// Namespace is the key to look for extra configuration details
const Namespace = "github_com/magalixtechnologies/krakend-gologging"

var (
	// ErrEmptyValue is the error returned when there is no config under the namespace
	ErrWrongConfig = fmt.Errorf("getting the extra config for the krakend-gologging module")
	// Used for krakend dependecies. Do not move
	DefaultPattern = ` %{time:2006/01/02 - 15:04:05.000} %{color}▶ %{level:.6s}%{color:reset} %{message}`
)

// NewLogger returns a krakend logger wrapping a gologging logger
func NewLogger(cfg config.ExtraConfig, ws ...io.Writer) (logging.Logger, error) {
	logConfig, ok := ConfigGetter(cfg).(Config)
	if !ok {
		return nil, ErrWrongConfig
	}
	logLevel := logConfig.Level
	var mgxlogLevel logger.Level
	switch logLevel {
	case "INFO":
		mgxlogLevel = logger.InfoLevel
	case "WARNING":
		mgxlogLevel = logger.WarnLevel
	case "DEBUG":
		mgxlogLevel = logger.DebugLevel
	case "ERROR":
		mgxlogLevel = logger.ErrorLevel
	default:
		return nil, fmt.Errorf("Unsupported log level %s", logLevel)
	}
	logger.Config(mgxlogLevel)
	return Logger{logLevel: mgxlogLevel}, nil
}

// ConfigGetter implements the config.ConfigGetter interface
func ConfigGetter(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return nil
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	cfg := Config{}
	if v, ok := tmp["level"]; ok {
		cfg.Level = v.(string)
	}
	return cfg
}

// Config is the custom config struct containing the params for the logger
type Config struct {
	Level string
}

// Logger is a wrapper over a github.com/op/go-logging logger
type Logger struct {
	logLevel logger.Level
}

// Debug implements the logger interface
func (l Logger) Debug(v ...interface{}) {
	logger.Debug(v...)
}

// Info implements the logger interface
func (l Logger) Info(v ...interface{}) {
	logger.Info(v...)
}

// Warning implements the logger interface
func (l Logger) Warning(v ...interface{}) {
	logger.Warn(v...)
}

// Error implements the logger interface
func (l Logger) Error(v ...interface{}) {
	logger.Error(v...)
}

// Critical implements the logger interface
func (l Logger) Critical(v ...interface{}) {
	logger.Error(v...)
}

// Fatal implements the logger interface
func (l Logger) Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

func (l Logger) GetLogLevel() logger.Level {
	return l.logLevel
}
