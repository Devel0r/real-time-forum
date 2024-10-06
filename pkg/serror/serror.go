package serror

import "errors"

// sentinel or signal errors
var (
	// Config errors
	ErrInvalidConfigPath = errors.New("error, invalid config path")
	ErrFileNotExists     = errors.New("error, a file is not exists")

	// Database errors
	ErrNilConfigStruct = errors.New("error, config structure is nil")
	ErrUserNotFound    = errors.New("error, user not found")
	ErrEmptyEnv        = errors.New("errror, empty ENV")

	// validation errors
	ErrIncorrectAge          = errors.New("error, user enter incorrect age")
	ErrIncorrectNameOrGender = errors.New("error, empty name or gender value")
	ErrInvalidEmail          = errors.New("error, invalid email")
	ErrEmptyEmail            = errors.New("error, empty email")
	ErrInvalidPassword       = errors.New("error, invalid password")  // errors.Is // Проверяет на схожесть ошибки
	ErrEmptyUsername         = errors.New("error, username is empty") // errors.As // Проверяет тип ошибки
	ErrEmptyFieldLogin       = errors.New("error, empty login or password user")

	// Save data errors
	ErrEmptyUserData   = errors.New("error, empty user data")
	ErrEmptyCookieData = errors.New("error, session structur pointer is nill")

	// Post errors
	ErrEmptyPostData           = errors.New("error, empty post data")
	ErrEmptyPostContentOrTitle = errors.New("error, empty post content or title")

	// Comment errors
	ErrEmptyCommentData = errors.New("error, empty comment data")

	// Categories errors
	ErrEmptyCategoryData = errors.New("error, empty category data")
)
