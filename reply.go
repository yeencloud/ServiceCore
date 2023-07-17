package servicecore

import (
	"github.com/gin-gonic/gin"
	"github.com/yeencloud/ServiceCore/domain"
	"github.com/yeencloud/ServiceCore/serviceError"
	"github.com/yeencloud/ServiceCore/tools"
)

func buildServiceReply() domain.ServiceReply {
	return domain.ServiceReply{
		ApiVersion: domain.APIVersion,
	}
}

func (shs *ServiceHTTPServer) replyWithError(c *gin.Context, requestId string, err *serviceError.Error, validationErrors []string) domain.ServiceReply {
	reply := buildServiceReply()
	reply.Error = err
	reply.ValidationErrors = tools.ArrayOrNil(validationErrors)

	c.AbortWithStatusJSON(err.HttpCode, reply)

	return reply
}