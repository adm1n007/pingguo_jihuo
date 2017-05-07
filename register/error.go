package register

import (
	. "ml/trace"
)

type GenericError struct {
	*BaseException
}

type RegisterError struct {
	*GenericError
}

type AlreadyExistsError struct {
	*GenericError
}

type PasswordTooSimpleError struct {
	*GenericError
}

type AccountBannedError struct {
	*GenericError
}

type AccountLockedSecurityError struct {
	*GenericError
}

type AccountDisabledError struct {
	*GenericError
}

type ExitingError struct {
	*GenericError
}

func NewGenericError(format string, args ...interface{}) *GenericError {
	return &GenericError{
		BaseException: NewBaseException(format, args...),
	}
}

func NewAlreadyExistsError(format string, args ...interface{}) *AlreadyExistsError {
	return &AlreadyExistsError{
		GenericError: NewGenericError(format, args...),
	}
}

func NewPasswordTooSimpleError(format string, args ...interface{}) *PasswordTooSimpleError {
	return &PasswordTooSimpleError{
		GenericError: NewGenericError(format, args...),
	}
}

func NewRegisterError(format string, args ...interface{}) *RegisterError {
	return &RegisterError{
		GenericError: NewGenericError(format, args...),
	}
}

func NewAccountBannedError(format string, args ...interface{}) *AccountBannedError {
	return &AccountBannedError{
		GenericError: NewGenericError(format, args...),
	}
}

func NewAccountLockedSecurityError(format string, args ...interface{}) *AccountLockedSecurityError {
	return &AccountLockedSecurityError{
		GenericError: NewGenericError(format, args...),
	}
}

func NewAccountDisabledError(format string, args ...interface{}) *AccountDisabledError {
	return &AccountDisabledError{
		GenericError: NewGenericError(format, args...),
	}
}

func NewExitingError() *ExitingError {
	return &ExitingError{
		GenericError: NewGenericError("Exiting"),
	}
}
