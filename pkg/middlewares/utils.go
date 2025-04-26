package middlewares

type Resp struct {
	ErrorCode   interface{} `json:"errorCode"`
	Message     string      `json:"message"`
	Description string      `json:"description"`
}

const (
	ErrDataInvalidV1       = "ErrDataInvalid"
	ErrIASAuthenticationV1 = "ErrIASAuthentication"
	ErrDataInvalid         = 1002

	ErrAuthentication = 4000

	ErrIASAuthentication = 4020

	ErrRegionNotFound = 4911
	ErrXTicketIdEmpty = 4912

	ErrCheckSig = 4921
)
