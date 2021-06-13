package gologging

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/MagalixTechnologies/core/logger"
	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/logging"
)

func getRequestId(r *http.Request) interface{} {
	reqID := r.Header.Get("X-Request-Id")
	if reqID == "" {
		reqID = shortID()
	}

	return reqID
}

func shortID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func NewRequestLogger(log logging.Logger) gin.HandlerFunc {
	lg := log.(Logger)
	level := lg.GetLogLevel()
	return func(c *gin.Context) {
		request := c.Request
		reqID := getRequestId(request)
		sugar := logger.New(lg.GetLogLevel())
		sugarLogger := sugar.With("requestId", reqID)
		started := time.Now()

		c.Next()

		if level == logger.DebugLevel {
			// copy request payload in case we will show it
			buf, _ := ioutil.ReadAll(request.Body)
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
			request.Body = rdr2 // OK since rdr2 implements the io.ReadCloser interface

			var payload map[string]interface{}
			json.NewDecoder(rdr1).Decode(&payload)

			headers := make(map[string][]string)
			for name, values := range request.Header {
				switch name {
				case "Accept", "X-Request-Id", "Content-Length", "Content-Type", "User-Agent", "X-B3-Parentspanid", "X-B3-Sampled", "X-B3-Spanid", "X-B3-Traceid", "X-Envoy-Attempt-Count", "X-Forwarded-Client-Cert", "X-Forwarded-Proto":
					continue
				default:
					headers[name] = values
				}
			}

			sugarLogger.Debugw("Default Log",
				"method", request.Method,
				"endpoint", request.URL.String(),
				"payload", payload,
				"headers", headers,
			)
		}

		defer sugarLogger.Sync()
		sugarLogger.Infow("Default Log",
			"method", request.Method,
			"endpoint", request.URL.String(),
			"StatusCode", c.Writer.Status(),
			"bytes", c.Writer.Size(),
			"duration", time.Since(started).String(),
		)
	}
}
