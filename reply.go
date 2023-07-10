package servicecore

import (
	"github.com/gin-gonic/gin"
	"github.com/yeencloud/ServiceCore/domain"
	"github.com/yeencloud/ServiceCore/tools"
)

func buildServiceReply() domain.ServiceReply {
	return domain.ServiceReply{
		Version: domain.APIVersion,
	}
}

func (shs *ServiceHTTPServer) replyWithError(c *gin.Context, err *domain.ServiceError, validationErrors []string) domain.ServiceReply {
	reply := buildServiceReply()
	reply.Error = err.Error.Error()
	reply.ValidationErrors = tools.ArrayOrNil(validationErrors)

	c.AbortWithStatusJSON(err.Code, reply)

	return reply
}