package xerrors

const (
	InvalidAdID  Error = "Ad ID is invalid"
	AdNotFound   Error = "Ad not found"
	MissingParam Error = "Missing param"
	AuthorizationError Error = "Authorization token is not in the headers"
	AdIsClose Error = "Ad is closed"
	Unauthorized Error = "You don't have the authorization"
	InternalServerError Error = "Intertal server error"
)

const (
	CodeInvalidData = iota
)
