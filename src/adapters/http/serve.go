package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	serviceError "github.com/yeencloud/ServiceCore/src/adapters/error"
	"github.com/yeencloud/ServiceCore/src/adapters/rpc"
	domain "github.com/yeencloud/ServiceCore/src/domain"
	errorDomain "github.com/yeencloud/ServiceCore/src/domain/serviceError"
	"github.com/yeencloud/ServiceCore/src/helpers"
	"net/http"
	"reflect"
)

// This middleware is used to fetch and fill the parameter for the method that will be called
func (shs *ServiceHTTPServer) getParameterStructFromBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Fetching the request struct from the body. This is the object that contains the parameters for the method that will be called.
		request, err := shs.getRequestStruct(c)
		if err != nil {
			return
		}

		//Fetching the service instance (the object - most likely a struct - that contains the methods that will be called)
		serviceInstance := reflect.ValueOf(shs.service)
		if !serviceInstance.IsValid() {
			return
		}

		//Fetching the method that will be called from the service instance by its name

		//We're getting its type
		methodType, _ := reflect.TypeOf(shs.service).MethodByName(request.Method)

		//and its value
		methodToCall := serviceInstance.MethodByName(request.Method)
		if !methodToCall.IsValid() {
			return
		}

		//Here we're filling the parameter struct with the values from the request struct and returning the validation errors if any arise.
		//validationErrors := rpc.FillParameterStruct(&r, input, request.Data)

		methodParameter, validationErrors := rpc.CreateMethodParameter(methodType.Type.In(1), request.Data)

		if len(validationErrors) > 0 {
			shs.replyWithError(c, request.RequestID, serviceError.Trace(domain.ErrValidationFailed), validationErrors)
			return
		}

		c.Set("parameter", methodParameter)
	}
}

func (shs *ServiceHTTPServer) checkRequestHasRequestStruct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request domain.ServiceRequest
		err := c.BindJSON(&request)

		if err != nil {
			shs.replyWithError(c, request.RequestID, serviceError.Trace(domain.ErrRequestCouldNotBindRequest), nil)
			return
		}

		c.Set("request", request)
	}
}

func CheckRequestVersion(request *domain.ServiceRequest) *errorDomain.Error {
	if request.ApiVersion != domain.APIVersion {
		return serviceError.Trace(domain.ErrVersionMismatch)
	}

	return nil
}

// this function should check for version
func (shs *ServiceHTTPServer) checkRequestIsValid() gin.HandlerFunc {
	return func(c *gin.Context) {
		request, err := shs.getRequestStruct(c)
		if err != nil {
			return
		}

		serr := CheckRequestVersion(request)

		if serr != nil {
			shs.replyWithError(c, request.RequestID, serr, nil)
		}

		err = shs.rpc.CheckRequestModule(request)

		if err != nil {
			shs.replyWithError(c, request.RequestID, err, nil)
		}

		service := request.Method

		methodFound := false
		for _, moduleService := range shs.serviceContent.Methods {
			if moduleService.Name == service {
				methodFound = true
				c.Set("requiredInput", moduleService.Input)
				break
			}
		}

		if !methodFound {
			shs.replyWithError(c, request.RequestID, serviceError.Trace(domain.ErrMethodNotFound), nil)
		}

		if request.Data == nil || len(request.Data) <= 0 {
			shs.replyWithError(c, request.RequestID, serviceError.Trace(domain.ErrRequestDataIsMissing), nil)
		}
	}
}

func (shs *ServiceHTTPServer) getRequestStruct(c *gin.Context) (*domain.ServiceRequest, *errorDomain.Error) {
	request, _ := c.Get("request")

	castRequest, succeeded := request.(domain.ServiceRequest)

	if !succeeded {
		shs.replyWithError(c, "", serviceError.Trace(domain.ErrRequestCouldNotBeCast), nil)
		return nil, serviceError.Trace(domain.ErrRequestCouldNotBeCast)
	}

	return &castRequest, nil
}

func (shs *ServiceHTTPServer) callServiceMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		//body request
		request, err := shs.getRequestStruct(c)
		if err != nil {
			return
		}

		if request.RequestID != "" {
			_, stderr := uuid.Parse(request.RequestID)
			if stderr != nil {
				shs.replyWithError(c, request.RequestID, serviceError.Trace(domain.ErrInvalidRequestID), nil)
				return
			}
		}

		serviceInstance := reflect.ValueOf(shs.service)
		if !serviceInstance.IsValid() {
			shs.replyWithError(c, request.RequestID, serviceError.Trace(domain.ErrServiceNotFound), nil)
			return
		}

		methodType, _ := reflect.TypeOf(shs.service).MethodByName(request.Method)
		methodToCall := serviceInstance.MethodByName(request.Method)
		_ = methodType
		if !methodToCall.IsValid() {
			shs.replyWithError(c, request.RequestID, serviceError.Trace(domain.ErrMethodNotFound), nil)
			return
		}

		inpass, found := c.Get("parameter")
		if !found || inpass == nil {
			//shs.replyWithError(c, &ErrCouldNotGetMethodParameter, nil)
			return
		}

		results := methodToCall.Call([]reflect.Value{reflect.ValueOf(inpass)})

		if len(results) == 2 {
			callResult := results[0]
			serr := results[1]

			reply := domain.ServiceReply{
				RequestID:  request.RequestID,
				Module:     request.Module,
				Method:     request.Method,
				ApiVersion: domain.APIVersion,
			}

			if serr.Type() == reflect.TypeOf(&errorDomain.Error{}) && !serr.IsNil() {
				err := serr.Interface().(*errorDomain.Error)

				log.Err(errors.New(err.String)).Msg("service method has errored")
				reply.Error = err

				request, found := c.Get("requestmetadata")

				if !found {
					request = domain.ServiceRequest{}
				}

				reply.Error.AdditionnalData = map[string]interface{}{
					"Parameters": inpass,
					"Request":    request,
				}
				c.IndentedJSON(err.HttpCode, reply)
			} else {
				data := callResult.Interface()
				rep := helpers.AnyToMap(data)
				reply.Data = rep
				c.IndentedJSON(http.StatusOK, reply)
			}
		}
	}
}
