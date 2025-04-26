package utils

const (
	PathHealth  = "/health"
	PathMetrics = "/metrics"
)

const (
	KeyTraceInfo               = "trace_info"
	HeaderXMeProfile           = "X-Me-Profile"
	KeyRateLimit               = "key-limit"
	KeyRequestBody             = "request_body"
	KeyResponseBody            = "response_body"
	JwtSub                     = "sub"
	JwtExp                     = "exp"
	KeyEchoContextRequestBody  = "echo_context_request_body"
	KeyEchoContextResponseBody = "echo_context_response_body"
	KeyMongoMultiConnName      = "mongo_multi_conn_name"
	KeyRegion                  = "X-Client-Region"
	KeySignature               = "signature"
	KeyXTicketId               = "X-Ticket-Id"
)

const (
	TagNameEncrypt = "encrypt"
	TagValEncrypt  = "true"
)

const (
	IASTypeService  = "SERVICE"
	IASTypeClient   = "CLIENT"
	IASTypePublic   = "PUBLIC"
	IASTypeInternal = "INTERNAL"
)

const (
	IASTokenExpireTypeLimited   = "LIMITED"
	IASTokenExpireTypeUnlimited = "UNLIMITED"
)

const (
	VGREncryptKey = "VGR_ENCRYPT_KEY"
)

const (
	XIASCode     = "X-IAS-Code"
	XRequestData = "X-Request-Data"
)
