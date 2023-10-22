package webserver

import (
	"fmt"
	"time"

	"github.com/Brodobrey-inc/TestService/config"
	"github.com/Brodobrey-inc/TestService/logging"
	"github.com/Brodobrey-inc/TestService/webserver/endpoints"
	"github.com/gin-gonic/gin"
)

func Initialize() *gin.Engine {
	switch config.ServiceConfig.DebugLevel {
	case "warning":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(ginLogger)

	setRoutes(r)
	return r
}

func StartServer(r *gin.Engine) {
	logging.LogInfo("Starting webserver", "host", "127.0.0.1", "port", config.ServiceConfig.ServicePort)
	if err := r.Run(fmt.Sprintf("%s:%d", "127.0.0.1", config.ServiceConfig.ServicePort)); err != nil {
		logging.LogError(err, "Failed to start webserver")
	}

	logging.LogFatalError(nil, "Webserver stopped")
}

func ginLogger(c *gin.Context) {
	// Start timer
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	// Process request
	c.Next()
	param := gin.LogFormatterParams{
		Request: c.Request,
		Keys:    c.Keys,
	}

	// Stop timer
	param.TimeStamp = time.Now()

	param.Path = path
	param.Latency = param.TimeStamp.Sub(start)
	param.ClientIP = c.ClientIP()
	param.Method = c.Request.Method
	param.StatusCode = c.Writer.Status()
	param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
	param.BodySize = c.Writer.Size()

	if raw != "" {
		path = path + "?" + raw
	}
	logging.LogDebug("WebServer Logging",
		"ClientIP", param.ClientIP,
		"method", param.Method,
		"path", param.Path,
		"protocol", param.Request.Proto,
		"statusCode", param.StatusCode,
		"latency", param.Latency,
		"responseSize", param.BodySize,
		"error", param.ErrorMessage)

	c.Next()
}

func setRoutes(r *gin.Engine) {
	r.POST("/create_person", endpoints.CreatePersonWebhook)
	r.GET("/list_persons", endpoints.ListPersonsWebhook)
	r.POST("/update_person/:person_uuid", endpoints.UpdatePersonWebhook)
	r.POST("/remove_person/:person_uuid", endpoints.RemovePersonWebhook)
}
