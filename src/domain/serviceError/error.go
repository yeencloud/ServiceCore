package serviceError

type ErrorDescription struct {
	HttpCode int
	String   string
}

type Error struct {
	ErrorDescription

	Stack           Stack
	Child           *Error                 `json:",omitempty"`
	AdditionnalData map[string]interface{} `json:",omitempty"`
	Fingerprint     string
}
