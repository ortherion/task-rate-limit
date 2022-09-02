package errors

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrTokenInvalid = errors.New("invalid token")

	ErrUserNotHavePermissions = errors.New("user does not have permissions")
	ErrUserNotSignatory       = errors.New("user not an signatory")
	ErrUserHasAlreadySigned   = errors.New("user has already signed task")
	ErrTaskIsAccepted         = errors.New("task is accepted")
	ErrTaskIsRejected         = errors.New("task is rejected")
	ErrCastUser               = errors.New("fail cast user")
	ErrCastTaskID             = errors.New("fail cast task id")
	ErrCastTopic              = errors.New("failed to cast topic")
	ErrCastEvent              = errors.New("failed to cast event")
	ErrSendGrpc               = errors.New("can`t send to grpc")
)
