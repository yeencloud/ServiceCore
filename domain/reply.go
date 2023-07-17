package domain

import (
	"github.com/yeencloud/ServiceCore/serviceError"
)

type ServiceReply struct {
	RequestID  string `json:",omitempty"`
	Module     string `json:",omitempty"`
	Method     string `json:",omitempty"`
	ApiVersion Version

	Error            *serviceError.Error `json:",omitempty"`
	ValidationErrors []string            `json:",omitempty"`

	Data map[string]interface{} `json:",omitempty"`
}