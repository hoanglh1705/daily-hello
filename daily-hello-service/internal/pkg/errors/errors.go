package errors

type AppError struct {
	Code    string
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

var (
	// General
	ErrInvalidInput = AppError{"INVALID_INPUT", "Invalid request data"}
	ErrUnauthorized = AppError{"UNAUTHORIZED", "Unauthorized"}
	ErrForbidden    = AppError{"FORBIDDEN", "You don't have permission to access this resource"}
	ErrNotFound     = AppError{"NOT_FOUND", "Resource not found"}
	ErrInternal     = AppError{"INTERNAL_ERROR", "Internal server error"}

	// Attendance
	ErrInvalidLocation  = AppError{"INVALID_LOCATION", "Not in allowed area"}
	ErrFakeGPS          = AppError{"FAKE_GPS", "Fake GPS detected"}
	ErrAlreadyCheckedIn = AppError{"ALREADY_CHECKED_IN", "Already checked in today"}
	ErrNotCheckedIn     = AppError{"NOT_CHECKED_IN", "You have not checked in yet"}

	// Auth
	ErrEmailExists    = AppError{"EMAIL_EXISTS", "Email already registered"}
	ErrInvalidCreds   = AppError{"INVALID_CREDENTIALS", "Invalid email or password"}
	ErrInvalidToken   = AppError{"INVALID_TOKEN", "Invalid or expired token"}

	// Branch
	ErrBranchNotFound = AppError{"BRANCH_NOT_FOUND", "Branch not found"}
)
