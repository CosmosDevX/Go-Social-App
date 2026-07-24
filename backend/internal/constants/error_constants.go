package constants

const (
	//auth errors
	InvalidPassword   = "INVALID_PASSWORD"
	TokenSignError    = "TOKEN_SIGN_ERROR"
	InvalidTokenError = "INVALID_TOKEN"
	AccessTokenError  = "ACCESS_TOKEN_ERROR"
	RefreshTokenError = "REFRESH_TOKEN_ERROR"

	//repository errors
	NotFound       = "NOT_FOUND"
	RequestTimeout = "REQUEST_TIMEOUT"
	SaveError      = "SAVE_ERROR"
	DeleteError    = "DELETE_ERROR"
	FindError      = "FIND_ERROR"
	UpdateError    = "UPDATE_ERROR"
	CreateError    = "CREATE_ERROR"

	//service errors
	ParseError = "PARSE_ERROR"
)
