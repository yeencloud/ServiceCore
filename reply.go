package ServiceCore

type ServiceReply struct {
	Module  string `json:",omitempty"`
	Service string `json:",omitempty"`
	Version ApiVersion

	Error            string   `json:",omitempty"`
	ValidationErrors []string `json:",omitempty"`

	Data map[string]interface{} `json:",omitempty"`
}