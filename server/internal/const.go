package internal

const (
	JSONRPCVersion = "2.0"
)

const (
	CodeInvalidParams = -32602
	CodeInternalError = -32603
)

const (
	MethodInitialize               = "initialize"
	MethodNotificationsInitialized = "notifications/initialized"
	MethodPing                     = "ping"
	MethodGetPrompt                = "prompts/get"
	MethodListPrompts              = "prompts/list"
)
