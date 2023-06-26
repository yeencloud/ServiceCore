package servicecore

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/yeencloud/ServiceCore/decompose"
	"net"
	"net/http"
	"time"
)

type ServiceHTTPServer struct {
	engine *gin.Engine

	service        any
	serviceContent *decompose.Module
}

func buildServiceReply() ServiceReply {
	return ServiceReply{
		Version: 1,
	}
}

func replyWithError(err error, validationErrors []string) ServiceReply {
	reply := buildServiceReply()
	reply.Error = err.Error()
	reply.ValidationErrors = validationErrors

	spew.Dump(reply)

	return reply
}

func (shs *ServiceHTTPServer) checkRequestHasRequestStruct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request ServiceRequest
		err := c.BindJSON(&request)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, replyWithError(ErrRequestIsMissing, nil))
			return
		}

		c.Set("request", request)
	}
}

func (shs *ServiceHTTPServer) getRequestStruct(c *gin.Context) (*ServiceRequest, error) {
	request, found := c.Get("request")
	if !found {
		return nil, ErrRequestIsMissing
	}

	castRequest, succeeded := request.(ServiceRequest)

	if !succeeded {
		return nil, ErrRequestMalformed
	}

	return &castRequest, nil
}

// this function should check for version
func (shs *ServiceHTTPServer) checkRequestVersionIsValid() gin.HandlerFunc {
	return func(c *gin.Context) {
		request, err := shs.getRequestStruct(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, replyWithError(err, nil))
			return
		}

		if request.ApiVersion != APIVersion {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, replyWithError(ErrVersionMismatch, nil))
		}
	}
}

func (shs *ServiceHTTPServer) checkRequestAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Todo: Check if request can be traced to a known microservice or a microservice with the same auth key
	}
}

func (shs *ServiceHTTPServer) checkRequestIfModuleExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		request, err := shs.getRequestStruct(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, replyWithError(err, nil))
			return
		}

		if request.Service != shs.serviceContent.Name {
			c.AbortWithStatusJSON(http.StatusNotFound, replyWithError(ErrorUnknownModule, nil))
		}
	}
}

func (shs *ServiceHTTPServer) checkRequestIfServiceExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		request, err := shs.getRequestStruct(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, replyWithError(err, nil))
			return
		}

		service := request.Method

		for _, moduleService := range shs.serviceContent.Methods {
			if moduleService.Name == service {
				c.Set("requiredInput", moduleService.Input)
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusNotFound, replyWithError(ErrUnknownMethod, nil))
	}
}

func (shs *ServiceHTTPServer) logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		code := c.Writer.Status()

		logWeight := log.Info()
		if code >= 300 {
			logWeight = log.Warn()
		}
		currentLog := logWeight.
			Str("path", path).
			Str("method", method).
			Str("ip", clientIP).
			Time("at", startTime).
			Int("code", code).
			TimeDiff("duration", endTime, startTime)

		currentLog.Msg("request served")
	}
}

func newServiceHttpServer(c *Config, service any, serviceContent *decompose.Module) *ServiceHTTPServer {
	server := ServiceHTTPServer{
		service:        service,
		serviceContent: serviceContent,
	}

	environment := c.GetEnvironment()

	switch environment {
	case EnvironmentDevelopment:
		gin.SetMode(gin.DebugMode)
	case EnvironmentProduction:
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	//
	r.Use(server.logRequest())
	rpc := r.Group("/rpc")

	rpc.Use(server.checkRequestHasRequestStruct())
	rpc.Use(server.checkRequestVersionIsValid())
	rpc.Use(server.checkRequestAuthentication())
	rpc.Use(server.checkRequestIfModuleExists())
	rpc.Use(server.checkRequestIfServiceExists())
	rpc.Use(server.getParameterStructFromBody())
	rpc.POST("/", server.callServiceMethod())

	server.engine = r

	return &server
}

func (s *ServiceHost) Listen() error {
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", s.Config.GetRPCPort()))

	fmt.Println("Listening on port", ln.Addr().String())

	err := http.Serve(ln, s.ServiceHttpServer.engine)
	if err != nil {
		return err
	}
	err = s.ServiceHttpServer.engine.Run()

	if err != nil {
		return err
	}

	return nil
}