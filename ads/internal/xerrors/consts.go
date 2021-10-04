package xerrors

const (
	InvalidAdID  Error = "Ad ID is invalid"
	AdNotFound   Error = "Ad not found"
	MissingParam Error = "Missing param"
	AuthorizationError Error = "Authorization token is not in the headers"
)

const (
	CodeInvalidData = iota
)
