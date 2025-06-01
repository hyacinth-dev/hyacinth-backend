package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrForbidden           = newError(403, "Forbidden")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	ErrEmailAlreadyUse          = newError(1001, "The email is already in use.")
	ErrUsernameAlreadyUse       = newError(1002, "The username is already in use.")
	ErrVnetTokenAlreadyUse      = newError(1003, "The vnet token is already in use.")
	ErrCannotDowngrade          = newError(1004, "Cannot downgrade to a lower plan.")
	ErrVnetLimitExceeded        = newError(1005, "Virtual network limit exceeded for your privilege level.")
	ErrVnetClientsLimitExceeded = newError(1006, "VNet clients limit would be exceeded after downgrade.")
	ErrOriginalPasswordNotMatch = newError(1007, "The original password does not match.")
	ErrUsernameConflict         = newError(1008, "Username is already in use by another user.")
)
