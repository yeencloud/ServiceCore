package servicecore

type ServiceReply struct {
	Module  string `json:",omitempty"`
	Service string `json:",omitempty"`
	Version Version

	Error            string   `json:",omitempty"`
	ValidationErrors []string `json:",omitempty"`

	Data map[string]interface{} `json:",omitempty"`
}