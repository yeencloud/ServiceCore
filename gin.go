package servicecore

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yeencloud/ServiceCore/config"
	"github.com/yeencloud/ServiceCore/decompose"
	"github.com/yeencloud/ServiceCore/rpc"
	"net"
	"net/http"
)

type ServiceHTTPServer struct {
	engine *gin.Engine

	service        any
	serviceContent *decompose.Module

	rpc *rpc.RPC
}

func newServiceHttpServer(c *config.Config, service any, serviceContent *decompose.Module) *ServiceHTTPServer {
	server := ServiceHTTPServer{
		service:        service,
		serviceContent: serviceContent,
	}

	environment := c.GetEnvironment()

	switch environment {
	case config.EnvironmentDevelopment:
		gin.SetMode(gin.DebugMode)
	case config.EnvironmentProduction:
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(server.logRequest())

	r.GET("/health", func(context *gin.Context) {
		//unimplemented
	})
	rpcRoute := r.Group("/rpc")
	rpcRoute.Use(server.checkRequestHasRequestStruct())
	rpcRoute.Use(server.checkRequestIsValid())
	rpcRoute.Use(server.getParameterStructFromBody())
	rpcRoute.POST("/", server.callServiceMethod())

	server.engine = r

	return &server
}

func (sh *ServiceHost) Listen() error {
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", sh.Config.GetRPCPort()))

	fmt.Println("Listening on port", ln.Addr().String())

	err := http.Serve(ln, sh.ServiceHttpServer.engine)
	if err != nil {
		return err
	}
	err = sh.ServiceHttpServer.engine.Run()

	if err != nil {
		return err
	}

	return nil
}