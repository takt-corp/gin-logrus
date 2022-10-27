// ginlogrus provides middleware for logging requests in the gin web framework. It supports, skipping
// routes that you don't want logged, such as health checks. Logrus can be configured globally.
package ginlogrus

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddlewareParams defines the configuration options for the logger middleware
type LoggerMiddlewareParams struct {
	// a list of paths to skip logging for
	SkipPaths []string
}

// LoggerMiddleware returns a middleware handler function that can be used with the gin
// router for logging requests.
func LoggerMiddleware(params LoggerMiddlewareParams) gin.HandlerFunc {
	// build a map of the skipped paths
	var skipMap map[string]struct{}

	if length := len(params.SkipPaths); length > 0 {
		skipMap = make(map[string]struct{}, length)

		for _, path := range params.SkipPaths {
			skipMap[path] = struct{}{}
		}
	}

	// return the handler function
	return func(c *gin.Context) {

		// we use the match path because it will match the value
		// defined on the router
		matchPath := c.FullPath()

		if _, ok := skipMap[matchPath]; !ok {
			// get the context of the request
			ctx := c.Request.Context()

			// get basic information about the request
			requestID := c.GetHeader("X-Request-ID")
			start := time.Now().UTC()
			path := c.Request.URL.Path
			hostname, err := os.Hostname()

			if err != nil {
				hostname = "unknown"
			}

			// builld the request logger
			logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
				"request_id": requestID,
				"hostname":   hostname,
			})

			// process request
			c.Next()

			// response metrics
			timestamp := time.Now().UTC()
			latency := timestamp.Sub(start)
			statusCode := c.Writer.Status()
			dataLength := c.Writer.Size()

			logger = logger.WithFields(logrus.Fields{
				"status_code": statusCode,
				"latency":     latency.Milliseconds(),
				"data_length": dataLength,
			})

			if len(c.Errors) > 0 {
				logger = logger.WithField("errors", c.Errors)
			}

			msg := fmt.Sprintf("[%s] %d %s (%dms)", timestamp.Format(time.RFC3339), statusCode, path, latency.Milliseconds())

			if statusCode >= 500 {
				logger.Error(msg)
			} else if statusCode >= 400 {
				logger.Warn(msg)
			} else {
				logger.Info(msg)
			}

		}
	}
}
