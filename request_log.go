package gologging

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
)

const (
	RequestLogNamespace  = "github_com/magalixtechnologies/gin-logger"
	RequestLogmoduleName = "gin-logger"
)

func NewRequestLogger(cfg config.ExtraConfig, logger logging.Logger, loggerConfig gin.LoggerConfig) gin.HandlerFunc {
	v, ok := RequestLogConfigGetter(cfg).(RequestLogConfig)
	if !ok {
		return gin.LoggerWithConfig(loggerConfig)
	}

	loggerConfig.SkipPaths = v.SkipPaths
	logger.Info(fmt.Sprintf("%s: total skip paths set: %d", RequestLogmoduleName, len(v.SkipPaths)))

	loggerConfig.Output = ioutil.Discard
	loggerConfig.Formatter = Formatter{logger.(Logger), v}.DefaultFormatter
	return gin.LoggerWithConfig(loggerConfig)
}

type Formatter struct {
	logger Logger
	config RequestLogConfig
}

func (f Formatter) DefaultFormatter(params gin.LogFormatterParams) string {
	logData := []interface{}{
		"method", params.Method,
		"endpoint", params.Path,
		"StatusCode", params.StatusCode,
		"duration", params.Latency,
	}

	f.logger.InfoWithArgs("Default logging", logData...)

	return ""
}

func RequestLogConfigGetter(e config.ExtraConfig) interface{} {
	v, ok := e[RequestLogNamespace]
	if !ok {
		return nil
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	cfg := defaultConfigGetter()
	if skipPaths, ok := tmp["skip_paths"].([]interface{}); ok {
		var paths []string
		for _, skipPath := range skipPaths {
			if path, ok := skipPath.(string); ok {
				paths = append(paths, path)
			}
		}
		cfg.SkipPaths = paths
	}

	return cfg
}

func defaultConfigGetter() RequestLogConfig {
	return RequestLogConfig{
		SkipPaths: []string{},
	}
}

type RequestLogConfig struct {
	SkipPaths []string
}
