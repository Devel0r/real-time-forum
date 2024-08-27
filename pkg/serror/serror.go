package serror

import "errors"

// sentinel or signal errors
var (
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
	ErrInvalidPassword       = errors.New("error, invalid password")  // errors.Is // Проверяет на схожесть ошибки
	ErrEmptyUsername         = errors.New("error, username is empty") // errors.As // Проверяет тип ошибки
)

// type Some struct {} // errors.As
// type Other struct {} // errors.As
// type Another struct {} // errors.As

// some := Some{} // errors.Is
