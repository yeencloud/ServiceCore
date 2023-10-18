package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	environment2 "github.com/yeencloud/ServiceCore/src/adapters/environment"
	"github.com/yeencloud/ServiceCore/src/adapters/reflect/decompose"
	"github.com/yeencloud/ServiceCore/src/adapters/rpc"
	"github.com/yeencloud/ServiceCore/src/config"
	"github.com/yeencloud/ServiceCore/src/domain/types"
	"net"
	"net/http"
)

type ServiceHTTPServer struct {
	engine *gin.Engine

	service        any
	serviceContent *decompose.Module

	rpc *rpc.RPC

	rpcConfig config.RpcServer
}

func NewServiceHttpServer(c config.RpcServer, service any, serviceContent *decompose.Module) *ServiceHTTPServer {
	server := ServiceHTTPServer{
		service:        service,
		serviceContent: serviceContent,
	}

	server.rpcConfig = c
	server.rpc = &rpc.RPC{
		Module: serviceContent.Name,
	}

	environment := environment2.GetEnvironment()

	switch environment {

	case types.EnvironmentDevelopment:
		gin.SetMode(gin.DebugMode)
	case types.EnvironmentProduction:
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

func (server *ServiceHTTPServer) Listen() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", int(server.rpcConfig.Port)))

	if err != nil {
		return err
	}

	log.Info().Str("Address", ln.Addr().String()).Msg("Now Listening !")

	err = http.Serve(ln, server.engine)
	if err != nil {
		return err
	}
	err = server.engine.Run()

	if err != nil {
		return err
	}

	return nil

}
