package http

import (
	"github.com/gin-gonic/gin"
	domain2 "github.com/yeencloud/ServiceCore/src/domain"
	"github.com/yeencloud/ServiceCore/src/domain/serviceError"
	"github.com/yeencloud/ServiceCore/src/helpers"
)

func buildServiceReply() domain2.ServiceReply {
	return domain2.ServiceReply{
		ApiVersion: domain2.APIVersion,
	}
}

func (shs *ServiceHTTPServer) replyWithError(c *gin.Context, requestId string, err *serviceError.Error, validationErrors []string) domain2.ServiceReply {
	reply := buildServiceReply()
	reply.Error = err
	reply.ValidationErrors = helpers.ArrayOrNil(validationErrors)

	c.AbortWithStatusJSON(err.HttpCode, reply)

	return reply
}
